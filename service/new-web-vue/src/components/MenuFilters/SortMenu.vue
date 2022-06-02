<template>
  <a href="#" class="navbar-filter__link" onclick="return false">
    <font-awesome-icon icon="sort" class="fa-fw" />
    <span class="navbar-filter__span">Sorting</span>
  </a>
  <ul class="navbar-filter-items">
    <li class="navbar-filter-item">
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
              query: { ...plat, ...liveplat, ...inac, ...sort },
            }"
            @click="changeFilter()"
            class="navbar-submenu-item__link"
            :class="{ active: !reg.reg }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              icon="earth-americas"
            />
            <span class="navbar-submenu-item__span">Name</span>
          </router-link>
        </li>

        <li class="navbar-submenu-item">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: { ...plat, ...liveplat, ...inac, ...sort },
            }"
            @click="changeFilter()"
            class="navbar-submenu-item__link"
            :class="{ active: !reg.reg }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              icon="earth-americas"
            />
            <span class="navbar-submenu-item__span">Youtube</span>
          </router-link>
        </li>

        <li class="navbar-submenu-item">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: { ...plat, ...liveplat, ...inac, ...sort },
            }"
            @click="changeFilter()"
            class="navbar-submenu-item__link"
            :class="{ active: !reg.reg }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              icon="earth-americas"
            />
            <span class="navbar-submenu-item__span">Youtube (Views)</span>
          </router-link>
        </li>

        <li class="navbar-submenu-item">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: { ...plat, ...liveplat, ...inac, ...sort },
            }"
            @click="changeFilter()"
            class="navbar-submenu-item__link"
            :class="{ active: !reg.reg }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              icon="earth-americas"
            />
            <span class="navbar-submenu-item__span">Twitch</span>
          </router-link>
        </li>

        <li class="navbar-submenu-item">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: { ...plat, ...liveplat, ...inac, ...sort },
            }"
            @click="changeFilter()"
            class="navbar-submenu-item__link"
            :class="{ active: !reg.reg }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              icon="earth-americas"
            />
            <span class="navbar-submenu-item__span">BiliBili</span>
          </router-link>
        </li>

        <li class="navbar-submenu-item">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: { ...plat, ...liveplat, ...inac, ...sort },
            }"
            @click="changeFilter()"
            class="navbar-submenu-item__link"
            :class="{ active: !reg.reg }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              icon="earth-americas"
            />
            <span class="navbar-submenu-item__span">BiliBili (Views)</span>
          </router-link>
        </li>

        <li class="navbar-submenu-item">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: { ...plat, ...liveplat, ...inac, ...sort },
            }"
            @click="changeFilter()"
            class="navbar-submenu-item__link"
            :class="{ active: !reg.reg }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              icon="earth-americas"
            />
            <span class="navbar-submenu-item__span">Twitter</span>
          </router-link>
        </li>
      </ul>
    </li>
    <li class="navbar-filter-item" v-if="platforms.length">
      <router-link
        :to="{
          params: { id: $route.params.id },
          query: { ...reg, ...plat, ...inac },
        }"
        @click="setSort()"
        class="navbar-filter-item__link"
        :class="{ active: !sort.sort }"
      >
        <font-awesome-icon
          class="fa-fw navbar-filter-item__svg"
          icon="arrow-down-a-z"
        />
        <span class="navbar-filter-item__span">Reverse Order</span>
      </router-link>
    </li>
    <li class="navbar-filter-item" v-if="platforms.length">
      <router-link
        :to="{
          params: { id: $route.params.id },
          query: { ...reg, ...plat, ...inac },
        }"
        @click="setSort()"
        class="navbar-filter-item__link"
        :class="{ active: !sort.sort }"
      >
        <font-awesome-icon class="fa-fw navbar-filter-item__svg" icon="ban" />
        <span class="navbar-filter-item__span">Disable Live First</span>
      </router-link>
    </li>
  </ul>
</template>

<script>
// Add arrow-down-a-z and arrow-up-z-a to font-awesome
import { library } from "@fortawesome/fontawesome-svg-core"
import {
  faArrowDownAZ,
  faArrowUpZA,
  faSort,
  faBan,
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
  faArrowUpZA,
  faYoutube,
  faTwitch,
  faBilibili,
  faTwitter,
  faSort,
  faBan
)

import { useMemberStore } from "@/stores/members.js"

export default {
  data() {
    return {
      reg: {},
      plat: {},
      liveplat: {},
      inac: {},
      sort: {},
    }
  },
  created() {
    this.$watch(
      () => this.$route.query,
      () => {
        this.reg = this.$route.query.reg ? { reg: this.$route.query.reg } : {}
        this.plat = this.$route.query.plat
          ? { plat: this.$route.query.plat }
          : {}
        this.liveplat = this.$route.query.liveplat
          ? { liveplat: this.$route.query.liveplat }
          : {}
        this.inac = this.$route.query.inac
          ? { inac: this.$route.query.inac }
          : {}
        this.sort = this.$route.query.sort
          ? { sort: this.$route.query.sort }
          : {}
      },
      { immediate: true }
    )
  },
  computed: {
    platforms() {
      return useMemberStore().members.config.menu.platform
    },
  },
  methods: {
    async setSort() {
      await new Promise((r) => setTimeout(r, 0))
      useMemberStore().sortingMembers()
    },
  },
}
</script>
