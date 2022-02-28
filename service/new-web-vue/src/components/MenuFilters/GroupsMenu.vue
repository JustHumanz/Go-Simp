<template>
  <a href="#" class="navbar-filter__link" onclick="return false">
    <font-awesome-icon
      class="fa-fw"
      icon="users"
      v-if="!current_group || !groupid"
    />
    <font-awesome-icon
      class="fa-fw"
      icon="user"
      v-else-if="current_group && groupid == 10"
    />
    <img
      :src="current_group.GroupIcon"
      :alt="current_group.GroupName"
      class="navbar-filter__img"
      loading="lazy"
      v-else-if="current_group && groupid && groupid != 10"
    />
    <span
      class="inline-block xs:hidden sm:!inline-block"
      v-if="!current_group || !groupid"
      >Groups</span
    >
    <span class="inline-block xs:hidden sm:!inline-block" v-else>{{
      (
        current_group.GroupName.charAt(0).toUpperCase() +
        current_group.GroupName.slice(1)
      ).replace("_", " ")
    }}</span>
  </a>
  <ul class="navbar-filter-items">
    <li class="navbar-filter-item">
      <a href="" class="navbar-filter-item__link">All Vtubers</a>
    </li>
  </ul>
</template>

<script>
import axios from "axios"
import Config from "../../config.json"
import { library } from "@fortawesome/fontawesome-svg-core"
import { faUsers, faUser } from "@fortawesome/free-solid-svg-icons"

library.add(faUsers, faUser)

// read props groupid
export default {
  props: {
    groupid: {
      type: String,
    },
  },
  data() {
    return {
      groups: [],
      current_group: null,
    }
  },
  async created() {
    await this.getGroupData()

    await this.groups.forEach((group) => {
      if (group.ID == this.groupid) {
        console.log(group)
        this.current_group = group
      }
    })

    console.log(this.groups)
    console.log(this.current_group)

    console.log(this.groupid)
  },
  methods: {
    async getGroupData() {
      if (this.groups.length > 0) return
      console.log("Fetching group data...")

      // this.cancelGroups = axios.CancelToken.source()

      const data_groups = await axios
        .get(Config.REST_API2 + "/groups/", {
          // cancelToken: this.cancelGroups.token,
        })
        .then((response) => response.data)
        .catch((error) => {
          if (!axios.isCancel(error)) this.error_msg = error.message
        })

      if (this.error_msg) return false

      // sort group data from GroupName
      data_groups.sort((a, b) => {
        if (a.GroupName < b.GroupName) return -1
        if (a.GroupName > b.GroupName) return 1
        return 0
      })

      // add "all vtubers" in the first position of "groups"
      data_groups.unshift({
        GroupName: "All Vtubers",
        GroupIcon: "",
      })

      this.groups = data_groups
      console.log(`Total group: ${this.groups.length}`)
    },
  },
}
</script>
