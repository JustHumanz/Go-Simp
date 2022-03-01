<template>
  <a
    href="#"
    class="navbar-filter__link"
    onclick="return false"
    @click="toggleMenu"
  >
    <font-awesome-icon
      class="fa-fw"
      icon="users"
      v-if="!current_group || !$route.params.id"
    />
    <font-awesome-icon
      class="fa-fw"
      icon="user"
      v-else-if="current_group && $route.params.id == 10"
    />
    <img
      :src="current_group.GroupIcon"
      :alt="current_group.GroupName"
      class="navbar-filter__img"
      loading="lazy"
      v-else-if="current_group && $route.params.id && $route.params.id != 10"
    />
    <span
      class="navbar-filter__span"
      v-if="!current_group || !$route.params.id"
      >Groups</span
    >
    <span class="navbar-filter__span" v-else>{{
      (
        current_group.GroupName.charAt(0).toUpperCase() +
        current_group.GroupName.slice(1)
      ).replace("_", " ")
    }}</span>
  </a>
  <ul class="navbar-filter-items peer-focus-within:scale-y-100">
    <li class="navbar-pending" v-if="groups.length < 1">
      <img
        :src="`/src/assets/loading/${Math.floor(Math.random() * 7)}.gif`"
        class="navbar-pending__img"
      />
    </li>
    <li class="navbar-filter-item" v-for="group in groups" :key="group.ID">
      <router-link
        :to="`/vtuber/${group.ID || ''}`"
        class="navbar-filter-item__link"
      >
        <font-awesome-icon
          class="fa-fw navbar-filter-item__svg"
          icon="users"
          v-if="!group.ID"
        />
        <font-awesome-icon
          class="fa-fw navbar-filter-item__svg"
          icon="user"
          v-else-if="group.ID && group.ID == 10"
        />
        <img
          v-else-if="group.ID && group.ID != 10"
          :src="group.GroupIcon"
          :alt="group.GroupName"
          class="navbar-filter-item__img"
          loading="lazy"
        />
        <span class="navbar-filter-item__span">
          {{
            (
              group.GroupName.charAt(0).toUpperCase() + group.GroupName.slice(1)
            ).replace("_", " ")
          }}
        </span>
      </router-link>
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
  data() {
    return {
      groups: [],
      current_group: null,
    }
  },
  async created() {
    await this.getGroupData()

    this.$watch(
      () => this.$route.params?.id,
      async () => {
        this.current_group = await this.groups.find(
          (group) => group.ID == this.$route.params.id
        )
        console.log(`Get Group: ${this.current_group.GroupName}`)
      },
      { immediate: true }
    )
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
        if (a.GroupName.toLowerCase() < b.GroupName.toLowerCase()) return -1
        if (a.GroupName.toLowerCase() > b.GroupName.toLowerCase()) return 1
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
