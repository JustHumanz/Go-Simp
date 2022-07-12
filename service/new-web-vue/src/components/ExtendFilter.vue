<template>
  <div class="extend-filter">
    <h3 class="extend-filter__title">Regions</h3>
    <div class="extend-filter_items">
      <div class="extend-filter__item" v-for="region in regions">
        <input
          type="checkbox"
          :id="`${region.code.toLowerCase()}_filter`"
          name="filter-check"
          :value="region.code.toLowerCase()"
          class="extend-filter__item-input"
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
    <h3 class="extend-filter__title">Platform</h3>
    <div class="extend-filter_items">
      <div class="extend-filter__item" v-if="platforms.includes('youtube')">
        <input
          type="checkbox"
          id="yt_plat"
          name="platform-check"
          value="youtube"
          class="extend-filter__item-input"
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
          value="twitch"
          class="extend-filter__item-input"
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
          value="bilibili"
          class="extend-filter__item-input"
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
        />
        <label for="select_reverse" class="extend-filter__item-label">
          <font-awesome-icon
            icon="circle-check"
            class="extend-filter__item-icon"
          />
          <span class="extend-filter__item-text">Select Except</span>
        </label>
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

export default {
  mounted() {
    console.log(regions)
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
  },
}
</script>

<style lang="scss" scoped>
.extend-filter {
  @apply fixed top-16 z-[15] mb-3 w-screen select-none bg-slate-100 py-2 px-5 shadow-md dark:bg-slate-400;

  &__title {
    @apply text-lg font-semibold text-white;
  }

  &__items {
    @apply flex items-center;
  }

  &__item {
    @apply m-0.5 inline-block;

    &-input {
      @apply hidden;

      &:checked ~ .extend-filter__item-label {
        @apply bg-white text-slate-700;
      }
    }

    &-label {
      @apply flex cursor-pointer items-center space-x-1 rounded-full border-2 border-white px-3 py-1 text-sm text-white transition-all duration-200 ease-in-out;
    }

    &-flag {
      @apply inline-block w-5 min-w-[1.25rem] object-contain drop-shadow-md;
    }
  }
}
</style>
