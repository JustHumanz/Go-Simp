<template>
  <a href="#" class="navbar-filter__link" onclick="return false">
    <font-awesome-icon icon="sort" class="fa-fw" />
    <span class="navbar-filter__span">Sorting</span>
  </a>
  <ul class="navbar-filter-items">
    <li class="navbar-filter-item" v-if="platforms.length || twitter">
      <a
        href="#"
        class="navbar-filter-item__link sub-menu"
        onclick="return false"
      >
        <font-awesome-icon class="fa-fw navbar-filter-item__svg" icon="globe" />
        <span class="navbar-filter-item__span">Order</span>
      </a>
      <ul class="navbar-submenu-items">
        <li class="navbar-submenu-item">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: urlSort(),
            }"
            @click="setSort()"
            class="navbar-submenu-item__link"
            :class="{ active: !sort || sort === 'name' }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              icon="arrow-down-a-z"
            />
            <span class="navbar-submenu-item__span">Name</span>
          </router-link>
        </li>

        <li class="navbar-submenu-item" v-if="platforms.includes('youtube')">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: urlSort(`yt`),
            }"
            @click="setSort()"
            class="navbar-submenu-item__link"
            :class="{ active: sort === 'yt' }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              :icon="[`fab`, `youtube`]"
            />
            <span class="navbar-submenu-item__span">Youtube</span>
          </router-link>
        </li>

        <li class="navbar-submenu-item" v-if="platforms.includes('youtube')">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: urlSort(`ytv`),
            }"
            @click="setSort()"
            class="navbar-submenu-item__link"
            :class="{ active: sort === 'ytv' }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              :icon="[`fab`, `youtube`]"
            />
            <span class="navbar-submenu-item__span">Youtube (Views)</span>
          </router-link>
        </li>

        <li class="navbar-submenu-item" v-if="platforms.includes('twitch')">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: urlSort(`tw`),
            }"
            @click="setSort()"
            class="navbar-submenu-item__link"
            :class="{ active: sort === 'tw' }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              :icon="[`fab`, `twitch`]"
            />
            <span class="navbar-submenu-item__span">Twitch</span>
          </router-link>
        </li>

        <li class="navbar-submenu-item" v-if="platforms.includes('bilibili')">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: urlSort(`bl`),
            }"
            @click="setSort()"
            class="navbar-submenu-item__link"
            :class="{ active: sort === 'bl' }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              :icon="[`fab`, `bilibili`]"
            />
            <span class="navbar-submenu-item__span">BiliBili</span>
          </router-link>
        </li>

        <li class="navbar-submenu-item" v-if="platforms.includes('bilibili')">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: urlSort(`blv`),
            }"
            @click="setSort()"
            class="navbar-submenu-item__link"
            :class="{ active: sort === 'blv' }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              :icon="[`fab`, `bilibili`]"
            />
            <span class="navbar-submenu-item__span">BiliBili (Views)</span>
          </router-link>
        </li>

        <li class="navbar-submenu-item" v-if="twitter">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: urlSort(`twr`),
            }"
            @click="setSort()"
            class="navbar-submenu-item__link"
            :class="{ active: sort === 'twr' }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              :icon="[`fab`, `twitter`]"
            />
            <span class="navbar-submenu-item__span">Twitter</span>
          </router-link>
        </li>
      </ul>
    </li>
    <li class="navbar-filter-item" v-if="platforms.length || twitter">
      <a
        href="#"
        class="navbar-filter-item__link sub-menu"
        onclick="return false"
      >
        <font-awesome-icon
          class="fa-fw navbar-filter-item__svg"
          icon="people-group"
        />
        <span class="navbar-filter-item__span">Group</span>
      </a>
      <ul class="navbar-submenu-items">
        <li class="navbar-submenu-item">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: toggleDisableLiveFirst(),
            }"
            @click="setSort()"
            class="navbar-filter-item__link"
            :class="{ active: live }"
          >
            <font-awesome-icon
              class="fa-fw navbar-filter-item__svg"
              icon="ban"
            />
            <span class="navbar-filter-item__span">Live First</span>
          </router-link>
        </li>
      </ul>
    </li>
    <li class="navbar-filter-item">
      <router-link
        :to="{
          params: { id: $route.params.id },
          query: reverseSorting(),
        }"
        @click="setSort()"
        class="navbar-filter-item__link"
        :class="{ active: reverse }"
      >
        <font-awesome-icon
          class="fa-fw navbar-filter-item__svg"
          icon="arrow-down-1-9"
        />
        <span class="navbar-filter-item__span">Reverse Order</span>
      </router-link>
    </li>
  </ul>
</template>

<script>
// Add arrow-down-a-z and arrow-up-z-a to font-awesome
import { library } from "@fortawesome/fontawesome-svg-core"
import {
  faArrowDownAZ,
  faArrowDown19,
  faArrowUp91,
  faSort,
  faBan,
  faPeopleGroup,
} from "@fortawesome/free-solid-svg-icons"
// Add fa brand for youtube, twitch, bilibili, twitter
import {
  faYoutube,
  faTwitch,
  faBilibili,
  faTwitter,
} from "@fortawesome/free-brands-svg-icons"

library.add(
  faArrowDownAZ,
  faArrowDown19,
  faArrowUp91,
  faYoutube,
  faTwitch,
  faBilibili,
  faTwitter,
  faSort,
  faBan,
  faPeopleGroup
)

import { useMemberStore } from "@/stores/members.js"

export default {
  created() {
    this.$watch(
      () => this.$route.query,
      () => {
        this.sort = this.$route.query.sort?.replace("-", " ") || null
        this.reverse = this.$route.query.sort?.includes("-") ?? false
        this.live = this.$route.query.live === undefined
      },
      { immediate: true }
    )
  },
  computed: {
    platforms() {
      return useMemberStore().sortMenu.platform
    },
    twitter() {
      return useMemberStore().sortMenu.twitter
    },
  },
  methods: {
    async setSort() {
      await new Promise((r) => setTimeout(r, 0))
      useMemberStore().sortingMembers()
    },
    urlSort(sorting = "") {
      const { reg, plat, liveplat, inac, sort, live } = this.$route.query

      // check when sort have "-" in first position
      sorting = sort?.match(/^-/) ? `-${sorting}` : sorting
      // if sorting = "-"
      sorting = sorting === "-undefined" ? "-name" : sorting

      const params = new Object()
      if (reg) params.reg = reg
      if (plat) params.plat = plat
      if (liveplat) params.liveplat = liveplat
      if (inac) params.inac = inac
      if ((!sort && sorting) || (sort && sorting !== "")) params.sort = sorting
      if (live) params.live = live

      return params
    },
    reverseSorting() {
      const { reg, plat, liveplat, inac, sort, live } = this.$route.query

      // check when sort have "-" in first position
      let sorting = sort?.match(/^-/) ? sort.replace("-", "") : `-${sort}`
      // if sorting = "-"
      sorting = sorting === "-undefined" ? "-name" : sorting
      // if sorting = "name"
      sorting = sorting === "name" ? "" : sorting

      const params = new Object()
      if (reg) params.reg = reg
      if (plat) params.plat = plat
      if (liveplat) params.liveplat = liveplat
      if (inac) params.inac = inac
      if ((!sort && sorting) || (sort && sorting !== "")) params.sort = sorting
      if (live) params.live = live

      return params
    },
    toggleDisableLiveFirst() {
      const { reg, plat, liveplat, inac, sort, live } = this.$route.query

      const params = new Object()
      if (reg) params.reg = reg
      if (plat) params.plat = plat
      if (liveplat) params.liveplat = liveplat
      if (inac) params.inac = inac
      if (sort) params.sort = sort
      if (!live) params.live = "false"

      return params
    },
  },
}
</script>
