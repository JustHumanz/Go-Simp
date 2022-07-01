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
          <a
            href="#"
            @click="setSort('name')"
            class="navbar-submenu-item__link"
            :class="{
              active: sorting.type === 'name',
            }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              icon="arrow-down-a-z"
            />
            <span class="navbar-submenu-item__span">Name</span>
          </a>
        </li>

        <li class="navbar-submenu-item" v-if="platforms.includes('youtube')">
          <a
            href="#"
            @click="setSort('youtube')"
            class="navbar-submenu-item__link"
            :class="{ active: sorting.type === 'youtube' }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              :icon="[`fab`, `youtube`]"
            />
            <span class="navbar-submenu-item__span">Youtube</span>
          </a>
        </li>

        <li class="navbar-submenu-item" v-if="platforms.includes('youtube')">
          <a
            href="#"
            @click="setSort(`youtube-views`)"
            class="navbar-submenu-item__link"
            :class="{ active: sorting.type === 'youtube_views' }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              :icon="[`fab`, `youtube`]"
            />
            <span class="navbar-submenu-item__span">Youtube (Views)</span>
          </a>
        </li>

        <li class="navbar-submenu-item" v-if="platforms.includes('twitch')">
          <a
            href="#"
            @click="setSort('twitch')"
            class="navbar-submenu-item__link"
            :class="{ active: sorting.type === 'twitch' }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              :icon="[`fab`, `twitch`]"
            />
            <span class="navbar-submenu-item__span">Twitch</span>
          </a>
        </li>

        <li class="navbar-submenu-item" v-if="platforms.includes('bilibili')">
          <a
            href="#"
            @click="setSort('bilibili')"
            class="navbar-submenu-item__link"
            :class="{ active: sorting.type === 'bilibili' }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              :icon="[`fab`, `bilibili`]"
            />
            <span class="navbar-submenu-item__span">BiliBili</span>
          </a>
        </li>

        <li class="navbar-submenu-item" v-if="platforms.includes('bilibili')">
          <a
            href="#"
            @click="setSort('bilibili_views')"
            class="navbar-submenu-item__link"
            :class="{ active: sorting.type === 'bilibili_views' }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              :icon="[`fab`, `bilibili`]"
            />
            <span class="navbar-submenu-item__span">BiliBili (Views)</span>
          </a>
        </li>

        <li class="navbar-submenu-item" v-if="twitter">
          <a
            href="#"
            @click="setSort('twitter')"
            class="navbar-submenu-item__link"
            :class="{ active: sorting.type === 'twitter' }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              :icon="[`fab`, `twitter`]"
            />
            <span class="navbar-submenu-item__span">Twitter</span>
          </a>
        </li>
      </ul>
    </li>
    <li class="navbar-filter-item">
      <a
        href="#"
        class="navbar-filter-item__link sub-menu"
        onclick="return false"
      >
        <font-awesome-icon
          class="fa-fw navbar-filter-item__svg"
          icon="people-group"
        />
        <span class="navbar-filter-item__span">Groups</span>
      </a>
      <ul class="navbar-submenu-items">
        <li class="navbar-submenu-item">
          <a
            href="#"
            @click="setSort('^live')"
            class="navbar-submenu-item__link"
            :class="{ active: sorting.live }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              icon="circle"
            />
            <span class="navbar-submenu-item__span">Live First</span>
          </a>
        </li>

        <li class="navbar-submenu-item">
          <a
            href="#"
            @click="setSort('^inactive')"
            class="navbar-submenu-item__link"
            :class="{ active: sorting.inactive }"
          >
            <font-awesome-icon
              class="fa-fw navbar-submenu-item__svg"
              icon="skull"
            />
            <span class="navbar-submenu-item__span">Inactive Last</span>
          </a>
        </li>
      </ul>
    </li>
    <li class="navbar-filter-item">
      <a
        href="#"
        @click="setSort('toggle')"
        class="navbar-filter-item__link"
        :class="{ active: sorting.order === 'desc' }"
      >
        <font-awesome-icon
          class="fa-fw navbar-filter-item__svg"
          icon="arrow-down-1-9"
        />
        <span class="navbar-filter-item__span">Reverse Order</span>
      </a>
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
    this.sorting = useMemberStore().sorting
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
    async setSort(query) {
      if (
        query !== "toggle" &&
        this.sorting.order === "desc" &&
        !query.includes("^")
      )
        query = `-${query}`

      if (this.sorting.order === "desc" && query === "toggle")
        query = this.sorting.type

      if (this.sorting.order === "asc" && query === "toggle")
        query = `-${this.sorting.type}`

      useMemberStore().changeSort(query)

      await new Promise((r) => setTimeout(r, 0))
      useMemberStore().sortingMembers()
    },
  },
}
</script>
