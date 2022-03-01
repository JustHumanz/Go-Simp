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
      <a
        href="#"
        class="navbar-filter-item__link sm:!w-64"
      >
        <font-awesome-icon class="fa-fw navbar-filter-item__svg" icon="arrow-down-a-z" />
        <span class="navbar-filter-item__span">Name (A - Z)</span>
      </a>
    </li>
    <li class="navbar-filter-item" v-if="filters">
      <a
        href="#"
        class="navbar-filter-item__link sm:!w-64"
      >
        <font-awesome-icon class="fa-fw navbar-filter-item__svg" icon="arrow-up-z-a" />
        <span class="navbar-filter-item__span">Name (Z - A)</span>
      </a>
    </li>
    <li class="navbar-filter-item" v-if="filters && filters.youtube">
      <a
        href="#"
        class="navbar-filter-item__link sm:!w-64"
      >
        <font-awesome-icon class="fa-fw navbar-filter-item__svg" :icon="['fab', 'youtube']" />
        <span class="navbar-filter-item__span">YouTube Subs (Bigger)</span>
      </a>
    </li>
    <li class="navbar-filter-item" v-if="filters && filters.youtube">
      <a
        href="#"
        class="navbar-filter-item__link sm:!w-64"
      >
        <font-awesome-icon class="fa-fw navbar-filter-item__svg" :icon="['fab', 'youtube']" />
        <span class="navbar-filter-item__span">YouTube Subs (Smaller)</span>
      </a>
    </li>
    <li class="navbar-filter-item" v-if="filters && filters.twitch">
      <a
        href="#"
        class="navbar-filter-item__link sm:!w-64"
      >
        <font-awesome-icon class="fa-fw navbar-filter-item__svg" :icon="['fab', 'twitch']" />
        <span class="navbar-filter-item__span">Twitch Followers (Bigger)</span>
      </a>
    </li>
    <li class="navbar-filter-item" v-if="filters && filters.twitch">
      <a
        href="#"
        class="navbar-filter-item__link sm:!w-64"
      >
        <font-awesome-icon class="fa-fw navbar-filter-item__svg" :icon="['fab', 'twitch']" />
        <span class="navbar-filter-item__span">Twitch Followers (Smaller)</span>
      </a>
    </li>
    <li class="navbar-filter-item"  v-if="filters && filters.bilibili">
      <a
        href="#"
        class="navbar-filter-item__link sm:!w-64"
      >
        <font-awesome-icon class="fa-fw navbar-filter-item__svg" :icon="['fab', 'bilibili']" />
        <span class="navbar-filter-item__span">Bilibili Fans (Bigger)</span>
      </a>
    </li>
    <li class="navbar-filter-item" v-if="filters && filters.bilibili">
      <a
        href="#"
        class="navbar-filter-item__link sm:!w-64"
      >
        <font-awesome-icon class="fa-fw navbar-filter-item__svg" :icon="['fab', 'bilibili']" />
        <span class="navbar-filter-item__span">Bilibili Fans (Smaller)</span>
      </a>
    </li>
    <li class="navbar-filter-item"  v-if="filters && filters.twitter">
      <a
        href="#"
        class="navbar-filter-item__link sm:!w-64"
      >
        <font-awesome-icon class="fa-fw navbar-filter-item__svg" :icon="['fab', 'twitter']" />
        <span class="navbar-filter-item__span">Twitter Followers (Bigger)</span>
      </a>
    </li>
    <li class="navbar-filter-item"  v-if="filters && filters.twitter">
      <a
        href="#"
        class="navbar-filter-item__link sm:!w-64"
      >
        <font-awesome-icon class="fa-fw navbar-filter-item__svg" :icon="['fab', 'twitter']" />
        <span class="navbar-filter-item__span">Twitter Followers (Smaller)</span>
      </a>
    </li>
  </ul>
</template>

<script>
// Add arrow-down-a-z and arrow-up-z-a to font-awesome
import { library } from '@fortawesome/fontawesome-svg-core'
import { faArrowDownAZ, faArrowUpZA } from '@fortawesome/free-solid-svg-icons'
// Add fa brand for youtube, twitch, bilibili, twitter
import { faYoutube, faTwitch, faBilibili, faTwitter } from '@fortawesome/free-brands-svg-icons'

library.add(faArrowDownAZ, faArrowUpZA, faYoutube, faTwitch, faBilibili, faTwitter)

export default {
  props: {
    filters: {
      type: Object
    }
  },
}
</script>
