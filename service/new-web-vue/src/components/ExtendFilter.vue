<template>
  <div class="bg-filter">
    <div class="extend-filter">
      <h3 class="extend-filter__title" v-if="regions.length > 1">Regions</h3>
      <div class="extend-filter_items" v-if="regions.length > 1">
        <div class="extend-filter__item" v-for="region in regions">
          <input
            type="checkbox"
            :id="`${region.code.toLowerCase()}_filter`"
            name="filter-check"
            :value="region.code.toLowerCase()"
            @change="changeLang"
            class="extend-filter__item-input"
            :checked="old_regions.includes(region.code.toLowerCase())"
          />
          <label
            :for="`${region.code.toLowerCase()}_filter`"
            class="extend-filter__item-label"
          >
            <img
              draggable="false"
              class="extend-filter__item-flag"
              :src="`/assets/flags/${region.code.toLowerCase()}.svg`"
              :alt="region.code"
              onerror="this.src='/assets/flags/none.svg'"
            />
            <span class="extend-filter__item-text">{{ region.name }}</span>
          </label>
        </div>
      </div>
      <h3 class="extend-filter__title" v-if="platforms.length > 1">Platform</h3>
      <div class="extend-filter_items" v-if="platforms.length > 1">
        <div class="extend-filter__item" v-if="platforms.includes('youtube')">
          <input
            type="checkbox"
            id="yt_plat"
            name="platform-check"
            value="yt"
            class="extend-filter__item-input"
            @change="changePlatform"
            :checked="old_platforms.includes('yt')"
          />
          <label for="yt_plat" class="extend-filter__item-label">
            <font-awesome-icon
              :icon="['fab', 'youtube']"
              class="extend-filter__item-icon"
            />
            <span class="extend-filter__item-text">Youtube</span>
          </label>
        </div>

        <div class="extend-filter__item" v-if="platforms.includes('twitch')">
          <input
            type="checkbox"
            id="tw_plat"
            name="platform-check"
            value="tw"
            class="extend-filter__item-input"
            @change="changePlatform"
            :checked="old_platforms.includes('tw')"
          />
          <label for="tw_plat" class="extend-filter__item-label">
            <font-awesome-icon
              :icon="['fab', 'twitch']"
              class="extend-filter__item-icon"
            />
            <span class="extend-filter__item-text">Twitch</span>
          </label>
        </div>

        <div class="extend-filter__item" v-if="platforms.includes('bilibili')">
          <input
            type="checkbox"
            id="bili_plat"
            name="platform-check"
            value="bl"
            class="extend-filter__item-input"
            @change="changePlatform"
            :checked="old_platforms.includes('bl')"
          />
          <label for="bili_plat" class="extend-filter__item-label">
            <font-awesome-icon
              :icon="['fab', 'bilibili']"
              class="extend-filter__item-icon"
            />
            <span class="extend-filter__item-text">Bilibili</span>
          </label>
        </div>

        <div class="extend-filter__item" v-if="platforms.length > 1">
          <input
            type="checkbox"
            id="select_reverse"
            name="select_reverse"
            class="extend-filter__item-input"
            @change="checkSelectOnly"
            :checked="select_only"
          />
          <label for="select_reverse" class="extend-filter__item-label">
            <font-awesome-icon
              icon="circle-check"
              class="extend-filter__item-icon"
            />
            <span class="extend-filter__item-text">Select Only</span>
          </label>
        </div>
      </div>
      <h3 class="extend-filter__title" v-if="live.length > 1">Live Platform</h3>
      <div class="extend-filter_items" v-if="live.length > 1">
        <div class="extend-filter__item" v-if="live.includes('youtube')">
          <input
            type="checkbox"
            id="yt_live"
            name="live-check"
            value="yt"
            class="extend-filter__item-input"
            @change="changeLive"
            :checked="old_live.includes('yt')"
          />
          <label for="yt_live" class="extend-filter__item-label">
            <font-awesome-icon
              :icon="['fab', 'youtube']"
              class="extend-filter__item-icon"
            />
            <span class="extend-filter__item-text">Youtube</span>
          </label>
        </div>

        <div class="extend-filter__item" v-if="live.includes('twitch')">
          <input
            type="checkbox"
            id="tw_live"
            name="live-check"
            value="tw"
            class="extend-filter__item-input"
            @change="changeLive"
            :checked="old_live.includes('tw')"
          />
          <label for="tw_live" class="extend-filter__item-label">
            <font-awesome-icon
              :icon="['fab', 'twitch']"
              class="extend-filter__item-icon"
            />
            <span class="extend-filter__item-text">Twitch</span>
          </label>
        </div>

        <div class="extend-filter__item" v-if="live.includes('bilibili')">
          <input
            type="checkbox"
            id="bili_live"
            name="live-check"
            value="bl"
            class="extend-filter__item-input"
            @change="changeLive"
            :checked="old_live.includes('bl')"
          />
          <label for="bili_live" class="extend-filter__item-label">
            <font-awesome-icon
              :icon="['fab', 'bilibili']"
              class="extend-filter__item-icon"
            />
            <span class="extend-filter__item-text">Bilibili</span>
          </label>
        </div>
      </div>

      <h3 class="extend-filter__title" v-if="live.length > 0 || inactive_check">
        Groups
      </h3>
      <div class="extend-filter_items" v-if="live.length > 0 || inactive_check">
        <div class="extend-filter__item" v-if="live.length > 0">
          <input
            type="checkbox"
            id="only_live"
            name="only_live"
            value="true"
            class="extend-filter__item-input"
            :checked="live_only"
            @change="groupChange"
          />
          <label for="only_live" class="extend-filter__item-label">
            <font-awesome-icon icon="users" class="extend-filter__item-icon" />
            <span class="extend-filter__item-text">Live Only</span>
          </label>
        </div>

        <div class="extend-filter__item" v-if="inactive_check">
          <input
            type="checkbox"
            id="inactive"
            name="inactive"
            value="true"
            class="extend-filter__item-input"
            :checked="inactive"
            @change="groupChange"
          />
          <label for="inactive" class="extend-filter__item-label">
            <font-awesome-icon icon="skull" class="extend-filter__item-icon" />
            <span class="extend-filter__item-text">Inactive Only</span>
          </label>
        </div>
      </div>

      <div class="button-area">
        <button @click="submit" class="button-area__button">Filter</button>
        <button @click="cancel" class="button-area__button">Cancel</button>
      </div>
    </div>
  </div>
</template>

<script>
import { library } from "@fortawesome/fontawesome-svg-core"
import { faCircleCheck } from "@fortawesome/free-solid-svg-icons"
library.add(faCircleCheck)

import { useMemberStore } from "@/stores/members.js"
import regions from "@/regions.json"
import { onBeforeRouteLeave } from "vue-router"

export default {
  setup() {
    onBeforeRouteLeave((to, from) => {
      document.body.style.overflow = "auto"
    })
  },
  data() {
    return {
      selected_regions: [],
      selected_platforms: [],
      selected_live: [],
      select_only: false,
      live_only: false,
      inactive: false,
    }
  },
  mounted() {
    const filters = useMemberStore().filter

    // make body overflow hidden
    document.body.style.overflow = "hidden"

    this.selected_regions = [...this.old_regions] || []
    this.selected_platforms = [...this.old_platforms] || []
    this.selected_live = [...this.old_live] || []
    this.select_only = filters.platform?.includes("-")
    this.live_only = this.$route.query.liveonly === "true"
    this.inactive = this.$route.query.inac === "true"
  },
  computed: {
    regions() {
      const filtersMenu = useMemberStore().menuFilter
      return filtersMenu.region.map((region) => {
        for (const r of regions) {
          if (r.code.toLowerCase() === region.toLowerCase()) return r
        }
      })
    },
    platforms() {
      const filtersMenu = useMemberStore().menuFilter
      return filtersMenu.platform
    },
    live() {
      const filtersMenu = useMemberStore().menuFilter
      return filtersMenu.live
    },
    old_regions() {
      const filters = useMemberStore().filter
      return filters.region || []
    },
    old_platforms() {
      const filters = useMemberStore().filter
      return filters.platform?.replace("-", "")?.split(",") || []
    },
    old_live() {
      const filters = useMemberStore().filter
      return filters.live?.split(",") || []
    },
    inactive_check() {
      return useMemberStore().menuFilter.inactive
    },
  },
  methods: {
    cancel() {
      useMemberStore().toggleadvanced()
      document.body.style.overflow = "auto"
    },
    changeLang(e) {
      const el = e.target

      const value = el.value
      const checked = el.checked || false

      if (checked) this.selected_regions.push(value)
      else this.selected_regions.splice(this.selected_regions.indexOf(value), 1)
    },
    changePlatform(e) {
      const el = e.target

      const value = el.value
      const checked = el.checked || false

      const platforms = document.querySelectorAll(
        "input[name='platform-check']"
      )

      if (!this.select_only) {
        for (const platform of platforms) {
          if (platform.value !== value) {
            platform.checked = false
            this.selected_platforms.splice(
              this.platforms.indexOf(platform.value),
              1
            )
          }
        }
      }

      if (checked) this.selected_platforms.push(value)
      else
        this.selected_platforms.splice(
          this.selected_platforms.indexOf(value),
          1
        )
    },
    checkSelectOnly(e) {
      this.select_only = e.target.checked || false

      const platforms = document.querySelectorAll(
        "input[name='platform-check']"
      )

      const check_ones =
        [...platforms].filter((platform) => platform.checked).length < 2

      if (!check_ones) {
        platforms.forEach((platform, index) => {
          // check index last
          if (index !== platforms.length - 1) {
            platform.checked = false
            this.selected_platforms.splice(
              this.platforms.indexOf(platform.value),
              1
            )
          }
        })
      }
    },
    changeLive(e) {
      const el = e.target

      const value = el.value
      const checked = el.checked || false

      if (checked) this.selected_live.push(value)
      else this.selected_live.splice(this.selected_live.indexOf(value), 1)
    },
    groupChange(e) {
      const liveOnly_el = document.querySelector("#only_live")
      const inactive_el = document.querySelector("#inactive")

      if (e.target.id === liveOnly_el.id) {
        this.live_only = e.target.checked || false
        inactive_el.checked = false
        this.inactive = false
      } else if (e.target.id === inactive_el.id) {
        this.inactive = e.target.checked || false
        liveOnly_el.checked = false
        this.live_only = false
      }
    },
    async submit() {
      const store = useMemberStore()

      store.toggleadvanced()
      document.body.style.overflow = "auto"

      const query = new Object()

      if (
        this.selected_regions.length > 0 &&
        this.selected_regions.length !== this.regions.length
      )
        query.reg = this.selected_regions.join(",")

      if (
        this.selected_platforms.length > 0 &&
        this.selected_platforms.length !== this.platforms.length
      )
        query.plat = this.select_only
          ? `-${this.selected_platforms[0]}`
          : this.selected_platforms.join(",")

      if (
        this.selected_live.length > 0 &&
        this.selected_live.length !== this.live.length
      )
        query.liveplat = this.selected_live.join(",")

      if (this.live_only) query.liveonly = "true"
      if (this.inactive) query.inac = "true"

      // route push
      this.$router.push({
        path: "/vtubers",
        params: { id: this.$route.params.id },
        query,
      })

      await new Promise((resolve) => setTimeout(resolve, 0))
      store.filterMembers()
      store.sortingMembers()
    },
  },
}
</script>

<style lang="scss" scoped>
.bg-filter {
  @apply fixed z-[15] h-screen w-screen bg-black/50;
}

.extend-filter {
  @apply fixed top-16 mb-3 w-screen select-none overflow-y-auto overflow-x-hidden bg-slate-100 py-2 px-5 shadow-md dark:bg-slate-500;
  max-height: calc(100vh - 64px);

  &__title {
    @apply text-lg font-semibold dark:text-white;

    &:not(:first-child) {
      @apply mt-2;
    }
  }

  &__items {
    @apply flex items-center;
  }

  &__item {
    @apply m-0.5 inline-block;

    &-input {
      @apply hidden;

      &:checked ~ .extend-filter__item-label {
        @apply bg-slate-500 text-white dark:bg-white dark:text-slate-700;
      }
    }

    &-label {
      @apply flex cursor-pointer items-center space-x-1 rounded-full border-2 border-slate-500 px-3 py-1 text-sm transition-all duration-200 ease-in-out dark:border-white dark:text-white;
    }

    &-flag {
      @apply inline-block w-5 min-w-[1.25rem] object-contain drop-shadow-md;
    }
  }
}

.button-area {
  @apply mt-3 flex justify-end;

  &__button {
    @apply mx-1 rounded-full bg-blue-400 px-3 py-1 text-sm text-white shadow-md transition-all duration-200 ease-in-out hover:translate-y-0.5 hover:shadow-sm dark:bg-white dark:text-slate-700;
  }
}
</style>
