<template>
  <a href="#" class="navbar-filter__link" onclick="return false">Sort</a>
  <ul class="navbar-filter-items">
    <li class="navbar-pending" v-if="!filters">
      <img
        :src="`/src/assets/loading/${Math.floor(Math.random() * 7)}.gif`"
        class="navbar-pending__img"
      />
    </li>
    <li class="navbar-filter-item" v-if="filters">
      <router-link
        :to="{
          params: { id: $route.params.id },
          query: { ...reg, ...plat, ...inac },
        }"
        class="navbar-filter-item__link sm:!w-64"
        :class="{ active: !sort.sort }"
      >
        <font-awesome-icon
          class="fa-fw navbar-filter-item__svg"
          icon="arrow-down-a-z"
        />
        <span class="navbar-filter-item__span">Name (A - Z)</span>
      </router-link>
    </li>
    <li class="navbar-filter-item" v-if="filters">
      <router-link
        :to="{
          params: { id: $route.params.id },
          query: { ...reg, ...plat, ...inac, sort: '-name' },
        }"
        class="navbar-filter-item__link sm:!w-64"
        :class="{ active: sort.sort == '-name' }"
      >
        <font-awesome-icon
          class="fa-fw navbar-filter-item__svg"
          icon="arrow-up-z-a"
        />
        <span class="navbar-filter-item__span">Name (Z - A)</span>
      </router-link>
    </li>
    <li class="navbar-filter-item" v-if="filters && filters.youtube">
      <router-link
        :to="{
          params: { id: $route.params.id },
          query: { ...reg, ...plat, ...inac, sort: 'yt' },
        }"
        class="navbar-filter-item__link sm:!w-64"
        :class="{ active: sort.sort == 'yt' }"
      >
        <font-awesome-icon
          class="fa-fw navbar-filter-item__svg"
          :icon="['fab', 'youtube']"
        />
        <span class="navbar-filter-item__span">YouTube Subs (Bigger)</span>
      </router-link>
    </li>
    <li class="navbar-filter-item" v-if="filters && filters.youtube">
      <router-link
        :to="{
          params: { id: $route.params.id },
          query: { ...reg, ...plat, ...inac, sort: '-yt' },
        }"
        class="navbar-filter-item__link sm:!w-64"
        :class="{ active: sort.sort == '-yt' }"
      >
        <font-awesome-icon
          class="fa-fw navbar-filter-item__svg"
          :icon="['fab', 'youtube']"
        />
        <span class="navbar-filter-item__span">YouTube Subs (Smaller)</span>
      </router-link>
    </li>
    <li class="navbar-filter-item" v-if="filters && filters.twitch">
      <router-link
        :to="{
          params: { id: $route.params.id },
          query: { ...reg, ...plat, ...inac, sort: 'tw' },
        }"
        class="navbar-filter-item__link sm:!w-64"
        :class="{ active: sort.sort == 'tw' }"
      >
        <font-awesome-icon
          class="fa-fw navbar-filter-item__svg"
          :icon="['fab', 'twitch']"
        />
        <span class="navbar-filter-item__span">Twitch Followers (Bigger)</span>
      </router-link>
    </li>
    <li class="navbar-filter-item" v-if="filters && filters.twitch">
      <router-link
        :to="{
          params: { id: $route.params.id },
          query: { ...reg, ...plat, ...inac, sort: '-tw' },
        }"
        class="navbar-filter-item__link sm:!w-64"
        :class="{ active: sort.sort == '-tw' }"
      >
        <font-awesome-icon
          class="fa-fw navbar-filter-item__svg"
          :icon="['fab', 'twitch']"
        />
        <span class="navbar-filter-item__span">Twitch Followers (Smaller)</span>
      </router-link>
    </li>
    <li class="navbar-filter-item" v-if="filters && filters.bilibili">
      <router-link
        :to="{
          params: { id: $route.params.id },
          query: { ...reg, ...plat, ...inac, sort: 'bl' },
        }"
        class="navbar-filter-item__link sm:!w-64"
        :class="{ active: sort.sort == 'bl' }"
      >
        <font-awesome-icon
          class="fa-fw navbar-filter-item__svg"
          :icon="['fab', 'bilibili']"
        />
        <span class="navbar-filter-item__span">Bilibili Fans (Bigger)</span>
      </router-link>
    </li>
    <li class="navbar-filter-item" v-if="filters && filters.bilibili">
      <router-link
        :to="{
          params: { id: $route.params.id },
          query: { ...reg, ...plat, ...inac, sort: '-bl' },
        }"
        class="navbar-filter-item__link sm:!w-64"
        :class="{ active: sort.sort == '-bl' }"
      >
        <font-awesome-icon
          class="fa-fw navbar-filter-item__svg"
          :icon="['fab', 'bilibili']"
        />
        <span class="navbar-filter-item__span">Bilibili Fans (Smaller)</span>
      </router-link>
    </li>
    <li class="navbar-filter-item" v-if="filters && filters.twitter">
      <router-link
        :to="{
          params: { id: $route.params.id },
          query: { ...reg, ...plat, ...inac, sort: 'twr' },
        }"
        class="navbar-filter-item__link sm:!w-64"
        :class="{ active: sort.sort == 'twr' }"
      >
        <font-awesome-icon
          class="fa-fw navbar-filter-item__svg"
          :icon="['fab', 'twitter']"
        />
        <span class="navbar-filter-item__span">Twitter Followers (Bigger)</span>
      </router-link>
    </li>
    <li class="navbar-filter-item" v-if="filters && filters.twitter">
      <router-link
        :to="{
          params: { id: $route.params.id },
          query: { ...reg, ...plat, ...inac, sort: '-twr' },
        }"
        class="navbar-filter-item__link sm:!w-64"
        :class="{ active: sort.sort == '-twr' }"
      >
        <font-awesome-icon
          class="fa-fw navbar-filter-item__svg"
          :icon="['fab', 'twitter']"
        />
        <span class="navbar-filter-item__span"
          >Twitter Followers (Smaller)</span
        >
      </router-link>
    </li>
  </ul>
</template>

<script>
// Add arrow-down-a-z and arrow-up-z-a to font-awesome
import { library } from "@fortawesome/fontawesome-svg-core"
import { faArrowDownAZ, faArrowUpZA } from "@fortawesome/free-solid-svg-icons"
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
  faTwitter
)

export default {
  props: {
    filters: {
      type: Object,
    },
  },
  data() {
    return {
      reg: {},
      plat: {},
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
}
</script>
