<template>
  <a href="#" class="navbar-filter__link" onclick="return false"> Filters </a>
  <ul class="navbar-filter-items">
    <li class="navbar-pending" v-if="!filters">
      <img
        :src="`/src/assets/loading/${Math.floor(Math.random() * 7)}.gif`"
        class="navbar-pending__img"
      />
    </li>
    <li class="navbar-filter-item" v-if="filters && filters.region.length > 1">
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
              query: { ...plat, ...inac, ...sort },
            }"
            class="navbar-submenu-item__link"
            :class="{ active: !reg.reg }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              icon="earth-americas"
            />
            <span class="navbar-submenu-item__span">All Regions</span>
          </router-link>
        </li>
        <li
          class="navbar-submenu-item"
          v-for="region in filters.region"
          :key="region"
        >
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: { reg: region.code, ...plat, ...inac, ...sort },
            }"
            class="navbar-submenu-item__link"
            :class="{ active: reg.reg == region.code }"
          >
            <img
              :src="`/src/assets/flags/${region.flagCode}.svg`"
              :alt="region.name"
              class="navbar-submenu-item__img"
            />
            <span class="navbar-submenu-item__span">{{ region.name }}</span>
          </router-link>
        </li>
      </ul>
    </li>
    <li class="navbar-filter-item" v-if="filters && platform.length > 1">
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
              query: { ...reg, ...inac, ...sort },
            }"
            class="navbar-submenu-item__link"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              icon="circle-play"
            />
            <span class="navbar-submenu-item__span">All Platform</span>
          </router-link>
        </li>
        <li class="navbar-submenu-item" v-if="filters.youtube">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: { ...reg, plat: 'yt', ...inac, ...sort },
            }"
            class="navbar-submenu-item__link"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              :icon="['fab', 'youtube']"
            />
            <span class="navbar-submenu-item__span">YouTube</span>
          </router-link>
        </li>
        <li class="navbar-submenu-item" v-if="filters.twitch">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: { ...reg, plat: 'tw', ...inac, ...sort },
            }"
            class="navbar-submenu-item__link"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              :icon="['fab', 'twitch']"
            />
            <span class="navbar-submenu-item__span">Twitch</span>
          </router-link>
        </li>
        <li class="navbar-submenu-item" v-if="filters.bilibili">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: { ...reg, plat: 'bl', ...inac, ...sort },
            }"
            class="navbar-submenu-item__link"
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
    <li class="navbar-filter-item" v-if="filters && filters.inactive">
      <a
        href="#"
        class="navbar-filter-item__link sub-menu"
        onclick="return false"
      >
        <font-awesome-icon
          class="fa-fw navbar-filter-item__svg"
          icon="filter"
        />
        <span class="navbar-filter-item__span">Other</span>
      </a>

      <ul class="navbar-submenu-items">
        <li class="navbar-submenu-item">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: { ...reg, ...plat, ...sort },
            }"
            class="navbar-submenu-item__link"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              icon="ban"
            />
            <span class="navbar-submenu-item__span">All</span>
          </router-link>
        </li>
        <li class="navbar-submenu-item">
          <router-link
            :to="{
              params: { id: $route.params.id },
              query: { ...reg, ...plat, inac: 'true', ...sort },
            }"
            class="navbar-submenu-item__link"
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
  faBan,
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
  faBan
)

export default {
  props: {
    filters: {
      type: Object,
      default: null,
    },
  },
  data() {
    return {
      platform: [],
      reg: {},
      plat: {},
      inac: {},
      sort: {},
    }
  },
  created() {
    this.$watch(
      () => this.filters,
      () => {
        this.platform = []

        if (this.filters) {
          if (this.filters.youtube) this.platform.push("youtube")
          if (this.filters.twitch) this.platform.push("twitch")
          if (this.filters.bilibili) this.platform.push("bilibili")
        }
      },
      { immediate: true }
    )

    this.$watch(
      () => this.$route.query,
      () => {
        this.reg = this.$route.query.reg ? { reg: this.$route.query.reg } : {}
        this.plat = this.$route.query.plat
          ? { plat: this.$route.query.plat }
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
  methods: {},
}
</script>
