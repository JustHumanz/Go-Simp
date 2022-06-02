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
      draggable="false"
      :src="current_group.GroupIcon"
      :alt="current_group.GroupName"
      class="navbar-filter__img"
      v-else-if="current_group && $route.params.id && $route.params.id != 10"
    />
    <span
      class="navbar-filter__span-mobile"
      v-if="!current_group || !$route.params.id"
      >Groups</span
    >
    <span class="navbar-filter__span-mobile" v-else>{{
      (
        current_group.GroupName.charAt(0).toUpperCase() +
        current_group.GroupName.slice(1)
      ).replace("_", " ")
    }}</span>
  </a>
  <ul class="navbar-filter-items peer-focus-within:scale-y-100">
    <li class="navbar-pending" v-if="groups.length < 1">
      <img
        draggable="false"
        :src="`/assets/loading/${Math.floor(Math.random() * 7)}.gif`"
        class="navbar-pending__img"
      />
    </li>
    <li class="navbar-filter-item" v-for="group in groups" :key="group.ID">
      <router-link
        :to="`/vtubers/${group.ID || ''}`"
        @click="getVtuberData(group.ID)"
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
          draggable="false"
          v-else-if="group.ID && group.ID != 10"
          :src="group.GroupIcon"
          :alt="group.GroupName"
          class="navbar-filter-item__img"
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
import { library } from "@fortawesome/fontawesome-svg-core"
import { faUsers, faUser } from "@fortawesome/free-solid-svg-icons"

library.add(faUsers, faUser)

import { useGroupStore } from "@/stores/groups"
import { useMemberStore } from "@/stores/members.js"

// read props groupid
export default {
  async created() {
    document.title = "List Vtubers - Vtbot"
  },
  computed: {
    groups() {
      return [
        { GroupName: "All Vtubers", GroupIcon: "" },
        ...useGroupStore().groups.data,
      ]
    },
    current_group() {
      const store = useGroupStore()

      const group =
        store.groups.data.find((group) => group.ID == this.$route.params?.id) ||
        null

      document.title = group
        ? this.convertToNormalText(group.GroupName) + " - List Vtubers"
        : "List Vtubers - Vtbot"

      if (group) console.log(`Get Group: ${group.GroupName}`)

      return group
    },
  },
  methods: {
    convertToNormalText(name) {
      return name.charAt(0).toUpperCase() + name.slice(1).replace("_", " ")
    },
    async getVtuberData(id) {
      await useMemberStore().fetchMembers(id || null)
      useMemberStore().filterMembers()
      useMemberStore().sortingMembers()
    },
  },
}
</script>
