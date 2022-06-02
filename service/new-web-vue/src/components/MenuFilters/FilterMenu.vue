<template>
  <a href="#" class="navbar-filter__link" onclick="return false">
    <font-awesome-icon class="fa-fw" icon="filter" />
    <span class="navbar-filter__span"> Filters</span>
  </a>
  <ul class="navbar-filter-items">
    <li class="navbar-filter-item" v-if="getRegions.length > 1">
      <a
        href="#"
        class="navbar-filter-item__link sub-menu"
        onclick="return false"
      >
        <font-awesome-icon class="fa-fw navbar-filter-item__svg" icon="globe" />
        <span class="navbar-filter-item__span">Region</span>
      </a>
      <ul class="navbar-submenu-items">
        <li class="navbar-submenu-item">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: urlParams({ region: '' }),
            }"
            @click="changeFilter()"
            class="navbar-submenu-item__link"
            :class="{ active: !reg }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              icon="earth-americas"
            />
            <span class="navbar-submenu-item__span">All Regions</span>
          </router-link>
        </li>
        <li class="navbar-submenu-item" v-for="region in getRegions">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: urlParams({ region: region.code }),
            }"
            @click="changeFilter()"
            class="navbar-submenu-item__link"
            :class="{ active: reg == region.code }"
          >
            <img
              draggable="false"
              :src="`/assets/flags/${region.flagCode}.svg`"
              :alt="region.name"
              class="navbar-submenu-item__img"
            />
            <span class="navbar-submenu-item__span">{{ region.name }}</span>
          </router-link>
        </li>
      </ul>
    </li>
    <li class="navbar-filter-item" v-if="platforms.length > 1">
      <a
        href="#"
        class="navbar-filter-item__link sub-menu"
        onclick="return false"
      >
        <font-awesome-icon class="fa-fw navbar-filter-item__svg" icon="video" />
        <span class="navbar-filter-item__span">Platform</span>
      </a>

      <ul class="navbar-submenu-items">
        <li class="navbar-submenu-item">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: urlParams({ platform: '' }),
            }"
            @click="changeFilter()"
            class="navbar-submenu-item__link"
            :class="{ active: !plat }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              icon="circle-play"
            />
            <span class="navbar-submenu-item__span">All Platform</span>
          </router-link>
        </li>
        <li class="navbar-submenu-item" v-if="platforms.includes(`youtube`)">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: urlParams({ platform: 'yt' }),
            }"
            @click="changeFilter()"
            class="navbar-submenu-item__link"
            :class="{ active: plat == 'yt' }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              :icon="['fab', 'youtube']"
            />
            <span class="navbar-submenu-item__span">YouTube</span>
          </router-link>
        </li>
        <li class="navbar-submenu-item" v-if="platforms.includes(`twitch`)">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: urlParams({ platform: 'tw' }),
            }"
            @click="changeFilter()"
            class="navbar-submenu-item__link"
            :class="{ active: plat == 'tw' }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              :icon="['fab', 'twitch']"
            />
            <span class="navbar-submenu-item__span">Twitch</span>
          </router-link>
        </li>
        <li class="navbar-submenu-item" v-if="platforms.includes(`bilibili`)">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: urlParams({ platform: 'bl' }),
            }"
            @click="changeFilter()"
            class="navbar-submenu-item__link"
            :class="{ active: plat == 'bl' }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              :icon="['fab', 'bilibili']"
            />
            <span class="navbar-submenu-item__span">Bilibili</span>
          </router-link>
        </li>
      </ul>
    </li>

    <li class="navbar-filter-item" v-if="livePlatforms.length > 0">
      <a
        href="#"
        class="navbar-filter-item__link sub-menu"
        onclick="return false"
      >
        <font-awesome-icon
          class="fa-fw navbar-filter-item__svg"
          icon="circle"
        />
        <span class="navbar-filter-item__span">Live Platform</span>
      </a>

      <ul class="navbar-submenu-items">
        <li class="navbar-submenu-item">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: urlParams({ live: '' }),
            }"
            @click="changeFilter()"
            class="navbar-submenu-item__link"
            :class="{ active: !liveplat }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              icon="users"
            />
            <span class="navbar-submenu-item__span">Show All</span>
          </router-link>
        </li>
        <li
          class="navbar-submenu-item"
          v-if="livePlatforms.includes(`youtube`)"
        >
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: urlParams({ live: '-yt,tw,bl' }),
            }"
            @click="changeFilter()"
            class="navbar-submenu-item__link"
            :class="{ active: liveplat == '-yt,tw,bl' }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              icon="video"
            />
            <span class="navbar-submenu-item__span">Show Live Only</span>
          </router-link>
        </li>
        <li
          class="navbar-submenu-item"
          v-if="livePlatforms.includes(`youtube`)"
        >
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: urlParams({ live: 'yt' }),
            }"
            @click="changeFilter()"
            class="navbar-submenu-item__link"
            :class="{ active: liveplat == 'yt' }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              :icon="['fab', 'youtube']"
            />
            <span class="navbar-submenu-item__span">YouTube</span>
          </router-link>
        </li>
        <li class="navbar-submenu-item" v-if="livePlatforms.includes(`twitch`)">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: urlParams({ live: 'tw' }),
            }"
            @click="changeFilter()"
            class="navbar-submenu-item__link"
            :class="{ active: liveplat == 'tw' }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              :icon="['fab', 'twitch']"
            />
            <span class="navbar-submenu-item__span">Twitch</span>
          </router-link>
        </li>
        <li
          class="navbar-submenu-item"
          v-if="livePlatforms.includes(`bilibili`)"
        >
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: urlParams({ live: 'bl' }),
            }"
            @click="changeFilter()"
            class="navbar-submenu-item__link"
            :class="{ active: liveplat == 'bl' }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              :icon="['fab', 'bilibili']"
            />
            <span class="navbar-submenu-item__span">Bilibili</span>
          </router-link>
        </li>
      </ul>
    </li>

    <li class="navbar-filter-item" v-if="inactiveCheck">
      <a
        href="#"
        class="navbar-filter-item__link sub-menu"
        onclick="return false"
      >
        <font-awesome-icon class="fa-fw navbar-filter-item__svg" icon="user" />
        <span class="navbar-filter-item__span">Activity Status</span>
      </a>

      <ul class="navbar-submenu-items">
        <li class="navbar-submenu-item">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: urlParams({ inactive: '' }),
            }"
            @click="changeFilter()"
            class="navbar-submenu-item__link"
            :class="{ active: !inac }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              icon="people-group"
            />
            <span class="navbar-submenu-item__span">All</span>
          </router-link>
        </li>
        <li class="navbar-submenu-item">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: urlParams({ inactive: 'false' }),
            }"
            @click="changeFilter()"
            class="navbar-submenu-item__link"
            :class="{ active: inac == 'false' }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              icon="user-check"
            />
            <span class="navbar-submenu-item__span">Active</span>
          </router-link>
        </li>
        <li class="navbar-submenu-item">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: urlParams({ inactive: 'true' }),
            }"
            @click="changeFilter()"
            class="navbar-submenu-item__link"
            :class="{ active: inac == 'true' }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              icon="skull"
            />
            <span class="navbar-submenu-item__span">Inactive</span>
          </router-link>
        </li>
      </ul>
    </li>
    <li class="navbar-filter-item">
      <a href="#" class="navbar-filter-item__link" onclick="return false">
        <font-awesome-icon
          class="fa-fw navbar-filter-item__svg"
          icon="plus-circle"
        />
        <span class="navbar-filter-item__span">Advanced</span>
      </a>
    </li>
    <li class="navbar-filter-item">
      <router-link
        router-link
        :to="{
          params: { id: $route.params.id },
          query: removeAll(),
        }"
        @click="changeFilter()"
        class="navbar-filter-item__link"
      >
        <font-awesome-icon
          class="fa-fw navbar-filter-item__svg"
          icon="arrows-rotate"
        />
        <span class="navbar-filter-item__span">Reset All Filters</span>
      </router-link>
    </li>
  </ul>
</template>

<script>
import { library } from "@fortawesome/fontawesome-svg-core"
import {
  faFilter,
  faGlobe,
  faVideo,
  faCircle,
  faSkull,
  faEarthAmericas,
  faCirclePlay,
  faPlusCircle,
  faArrowsRotate,
  faUser,
  faUserCheck,
  faPeopleGroup,
} from "@fortawesome/free-solid-svg-icons"

// Add icon youtube, twitch, and bilibili from font-awesome-brands
import {
  faYoutube,
  faTwitch,
  faBilibili,
} from "@fortawesome/free-brands-svg-icons"

library.add(
  faFilter,
  faGlobe,
  faVideo,
  faCircle,
  faSkull,
  faYoutube,
  faTwitch,
  faBilibili,
  faEarthAmericas,
  faCirclePlay,
  faPlusCircle,
  faArrowsRotate,
  faUser,
  faUserCheck,
  faPeopleGroup
)

import regionConfig from "@/region.json"
import { useMemberStore } from "@/stores/members.js"

export default {
  created() {
    this.$watch(
      () => this.$route.query,
      () => {
        // set all data this.$route.query to this
        this.reg = this.$route.query.reg
        this.plat = this.$route.query.plat
        this.liveplat = this.$route.query.liveplat
        this.inac = this.$route.query.inac
      },
      { immediate: true }
    )
  },
  computed: {
    getRegions() {
      const store = useMemberStore()

      const regions = store.members.config.menu.region.map(
        (region) =>
          regionConfig.find(
            (r) => r.code.toLowerCase() == region.toLowerCase()
          ) || {
            code: region,
            name: `${region} Region`,
          }
      )

      regions.sort((a, b) => (a.name < b.name ? -1 : 1))

      return regions
    },

    platforms() {
      return useMemberStore().members.config.menu.platform
    },
    livePlatforms() {
      return useMemberStore().members.config.menu.live
    },
    inactiveCheck() {
      return useMemberStore().members.config.menu.inactive
    },
  },
  methods: {
    async changeFilter() {
      const store = useMemberStore()

      await new Promise((resolve) => setTimeout(resolve, 0))
      store.filterMembers()
      store.sortingMembers()
    },
    urlParams({
      region = null,
      platform = null,
      live = null,
      inactive = null,
    }) {
      const { reg, plat, liveplat, inac, sort } = this.$route.query

      const params = new Object()

      if ((!reg && region) || (reg && region !== "")) params.reg = region || reg
      if ((!plat && platform) || (plat && platform !== ""))
        params.plat = platform || plat
      if ((!liveplat && live) || (liveplat && live !== ""))
        params.liveplat = live || liveplat
      if ((!inac && inactive) || (inac && inactive !== ""))
        params.inac = inactive || inac
      if (sort) params.sort = sort

      return params
    },
    removeAll() {
      const { sort } = this.$route.query

      const params = new Object()
      if (sort) params.sort = sort

      return params
    },
  },
}
</script>
