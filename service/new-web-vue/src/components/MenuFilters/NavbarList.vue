<script setup>
import GroupsMenu from "./GroupsMenu.vue"
import FilterMenu from "./FilterMenu.vue"
import SortMenu from "./SortMenu.vue"
</script>

<template>
  <nav class="list-nav">
    <ul class="navbar-filters">
      <li
        class="navbar-filter group"
        :class="{
          disabled:
            (disabled && error_status != 404) ||
            (error_status && error_status != 404),
        }"
      >
        <GroupsMenu :groups="groups" />
      </li>
      <li
        class="navbar-filter"
        :class="{
          disabled:
            disabled ||
            error_status ||
            (filters &&
              filters.region.length < 2 &&
              platform.length < 2 &&
              !filters.inactive),
        }"
      >
        <FilterMenu :filters="filters" />
      </li>
      <li class="navbar-filter" :class="{ disabled: disabled || error_status }">
        <SortMenu :filters="filters" />
      </li>
    </ul>
    <div class="nav-search">
      <font-awesome-icon
        icon="magnifying-glass"
        class="fa-fw fa-md nav-search__svg"
      />
      <input
        type="text"
        class="nav-search__input"
        @keydown="searchData"
        ref="search_input"
        :placeholder="placeholder || `Search Vtubers...`"
        :disabled="disable_search"
      />
    </div>
  </nav>
</template>

<script>
import { library } from "@fortawesome/fontawesome-svg-core"
import { faMagnifyingGlass } from "@fortawesome/free-solid-svg-icons"

library.add(faMagnifyingGlass)

export default {
  data() {
    return {
      platform: [],
      search_query: null,
      err_status: null,
    }
  },
  props: {
    groups: {
      type: Array,
      default: [],
    },
    filters: {
      type: Object,
      default: null,
    },
    placeholder: {
      type: String,
      default: "Search Vtubers...",
    },
    disable_search: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    error_status: {
      type: Number,
      default: null,
    },
  },
  async created() {
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

    // this.$watch(
    //   () => this.search_query,
    //   () => {
    //     window.scrollTo({
    //       top: 0,
    //     })
    //     this.$emit("search", this.search_query)
    //   },
    //   { immediate: true }
    // )

    this.$watch(
      () => this.$route,
      (a, b) => {
        if (
          a.query.reg !== b.query.reg ||
          a.query.plat !== b.query.plat ||
          a.query.inac !== b.query.inac ||
          a.params.id !== b.params.id
        ) {
          this.$refs.search_input.value = ""
          this.$emit("search", null)
        }
      }
    )
  },
  methods: {
    async searchData() {
      await new Promise((resolve) => setTimeout(resolve, 60))
      this.$emit("search", this.$refs.search_input.value)
    },
  },
}
</script>

<style lang="scss">
.list-nav {
  @apply fixed top-16 z-10 flex w-screen select-none flex-wrap-reverse items-center justify-center bg-blue-400 py-2 px-5 dark:bg-slate-500 sm:justify-around;
}

.navbar-filters {
  @apply flex items-center space-x-1 first:mt-2 sm:space-x-2 xs:first:mt-0;
}

.navbar-pending {
  @apply flex h-[100vw] w-screen items-center justify-center p-3 sm:h-[10.5rem] sm:w-[10.5rem];
}

.navbar-filter {
  @apply sm:relative;

  &__link {
    @apply flex items-center space-x-1 rounded-md px-2 py-1 font-semibold text-white transition-all duration-200 ease-in-out;
  }

  &__img {
    @apply inline-block w-6 min-w-[1.5rem] object-contain drop-shadow-md;
  }

  &__svg {
    @apply w-6;
  }

  &__span {
    @apply inline-block xs:hidden;

    @media (min-width: 640px) {
      display: inline-block !important;
    }
  }

  &.disabled {
    @apply opacity-50;

    &__link {
      @apply cursor-not-allowed;
    }
  }

  &:not(.disabled) {
    .navbar-filter__link {
      @apply shadow-blue-600/75 hover:-translate-y-px hover:shadow-sm dark:shadow-slate-300/50;
    }

    &:focus-within {
      .navbar-filter__link {
        @apply -translate-y-0.5 shadow-md shadow-blue-600/75 dark:shadow-slate-300/50;
      }

      .navbar-filter-items {
        @apply scale-y-100;
      }
    }
  }

  &-items {
    @apply absolute left-0 mt-2 flex max-h-[83.8vh] origin-top scale-y-0 flex-col overflow-y-auto overflow-x-hidden bg-blue-400 transition-all dark:bg-slate-700 sm:left-auto sm:mt-0 sm:max-h-60 sm:rounded-md sm:shadow-md sm:shadow-blue-600/75 sm:dark:shadow-slate-200/75;
    @media (min-width: 640px) {
      scrollbar-width: none; /* Firefox */
      -ms-overflow-style: none; /* IE 10+ */
      &::-webkit-scrollbar {
        /* Chromium and Safari */
        display: none;
      }
    }
  }

  // add class exept sort
  &.group {
    .router-link-active {
      @apply bg-blue-600 dark:bg-slate-900;
    }
  }

  &-item {
    @apply bg-blue-400 dark:bg-slate-700;
    &__img {
      @apply inline-block w-5 min-w-[1.25rem] rounded-md object-contain drop-shadow-md;
    }

    &__svg {
      @apply w-6;
    }

    &__link {
      @apply flex w-screen items-center space-x-2 px-2 py-1 font-semibold text-white hover:bg-blue-500/60 dark:hover:bg-slate-900/40 sm:w-44;

      &.active {
        @apply bg-blue-600 dark:bg-slate-900;
      }

      &.sub-menu::after {
        // add arrow right icon
        @apply absolute right-3 rotate-90 border-y-5 border-l-5 border-solid border-y-transparent border-l-current transition content-[''];
      }
    }
    &:focus-within {
      .sub-menu::after {
        @apply rotate-0;
      }
      .navbar-submenu-items {
        @apply h-[var(--totalHeight)] scale-y-100;
      }
    }
  }
}

.navbar-submenu {
  &-items {
    transition-property: transform height;
    @apply flex h-0 origin-top scale-y-0 flex-col duration-200 ease-in-out;
  }

  &-item {
    @apply flex w-full items-center bg-blue-500/20 dark:bg-slate-900/20;

    &__link {
      @apply flex w-full items-center space-x-2 px-2 py-1 pl-7 font-semibold text-white hover:bg-blue-500/40 dark:hover:bg-slate-900/40 sm:pl-4;

      &.active {
        @apply bg-blue-600 dark:bg-slate-900;
      }
    }

    &__img {
      @apply inline-block w-5 min-w-[1.25rem] rounded-sm object-contain;
    }

    &__svg {
      @apply w-6;
    }
  }
}

.nav-search {
  @apply relative mx-1 ml-3 inline-block flex-auto hover:-translate-y-px sm:flex-none;

  &:focus-within {
    transform: translate(0, -2px) !important;
  }

  &__svg {
    @apply absolute mt-2 ml-2 text-blue-500 dark:text-white;
  }

  &__input {
    @apply w-full rounded-lg bg-blue-300 py-1 px-2 pl-8 font-semibold text-gray-600 transition-all placeholder:font-normal placeholder:italic placeholder:text-blue-500 hover:shadow-sm hover:shadow-blue-600/75 focus:bg-blue-200 focus:shadow-md  focus:shadow-blue-600/75 focus:outline-none disabled:bg-blue-600 disabled:placeholder:text-blue-200 dark:bg-slate-300/20 dark:text-white dark:placeholder:text-gray-300 dark:hover:shadow-slate-100/75 dark:focus:bg-slate-400 dark:focus:shadow-slate-100/75 dark:disabled:bg-slate-700;
  }
}
</style>
