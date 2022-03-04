<script setup>
import GroupsMenu from "./GroupsMenu.vue"
import FilterMenu from "./FilterMenu.vue"
import SortMenu from "./SortMenu.vue"
</script>

<template>
  <nav class="list-nav">
    <ul class="navbar-filters">
      <li class="navbar-filter group">
        <GroupsMenu />
      </li>
      <li
        class="navbar-filter"
        v-if="
          !filters ||
          !(
            filters.region.length < 2 &&
            platform.length < 2 &&
            !filters.inactive
          )
        "
      >
        <FilterMenu :filters="filters" />
      </li>
      <li class="navbar-filter">
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
        v-model="search_query"
        :placeholder="placeholder"
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
      activeMenu: null,
      platform: [],
      search_query: "",
    }
  },
  props: {
    filters: {
      type: Object,
      default: null,
    },
    placeholder: {
      type: String,
      default: "Search Vtuber...",
    },
    disable_search: {
      type: Boolean,
      default: false,
    },
  },
  created() {
    this.getClickMenu()
    this.unfocusMenu()

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
      () => this.search_query,
      () => {
        this.$emit("search", this.search_query)
      },
      { immediate: true }
    )

    this.$watch(
      () => this.$route,
      (a, b) => {
        if (
          a.query.reg !== b.query.reg ||
          a.query.plat !== b.query.plat ||
          a.query.inac !== b.query.inac ||
          a.params.id !== b.params.id
        )
          this.search_query = ""
      }
    )
  },
  methods: {
    getClickMenu() {
      document.onclick = (e) => {
        let classList = [...e.target.classList]

        if (e.target.tagName === "path") {
          classList = [...e.target.parentElement.parentElement.classList]
        }
        if (e.target.tagName === "svg") {
          classList = [...e.target.parentElement.classList]
        }

        if (classList.find((c) => c.includes("navbar-filter__"))) {
          const navbarFilter =
            e.target.tagName === "A"
              ? e.target
              : e.target.tagName === "path"
              ? e.target.parentElement.parentElement
              : e.target.parentElement

          switch (this.activeMenu) {
            case navbarFilter:
              this.activeMenu.blur()
              this.activeMenu = null
              break
            case null:
              this.activeMenu = navbarFilter
              break
            default:
              this.activeMenu = navbarFilter
              break
          }
        } else if (classList.find((c) => c.includes("navbar-filter-item__"))) {
          const navbarFilterItem =
            e.target.tagName === "A"
              ? e.target
              : e.target.tagName === "path"
              ? e.target.parentElement.parentElement
              : e.target.parentElement

          if (!navbarFilterItem.classList.contains("sub-menu")) {
            this.activeMenu = null
            navbarFilterItem.blur()
          }
        } else if (classList.find((c) => c.includes("navbar-submenu-item__"))) {
          const navbarSubItem =
            e.target.tagName === "A"
              ? e.target
              : e.target.tagName === "path"
              ? e.target.parentElement.parentElement
              : e.target.parentElement

          this.activeMenu = null
          navbarSubItem.blur()
        } else if (classList.find((c) => c.includes("nav-search"))) {
          if (this.activeMenu !== null) {
            console.log("closing menu")
            this.activeMenu = null
          }

          const navbarSearchItem =
            e.target.tagName === "DIV"
              ? e.target
              : e.target.tagName === "path"
              ? e.target.parentElement.parentElement
              : e.target.parentElement

          navbarSearchItem.children[1].focus()
        } else {
          if (this.activeMenu === null) return
          console.log("closing menu")
          this.activeMenu = null
        }
      }
    },
    unfocusMenu() {
      // when document unfocus
      document.onblur = (e) => {
        if (this.activeMenu) {
          if (this.activeMenu === document.activeElement) this.activeMenu.blur()
          else if (
            !document.activeElement.classList.contains("nav-search__input")
          )
            document.activeElement.blur()

          console.log("closing menu")
          this.activeMenu = null
        }
      }
    },
  },
}
</script>

<style lang="scss">
.list-nav {
  @apply bg-blue-400 fixed top-16 py-2 px-5 w-screen flex flex-wrap-reverse items-center sm:justify-around justify-center z-10 select-none;
}

.navbar-filters {
  @apply flex space-x-1 sm:space-x-2 first:mt-2 xs:first:mt-0 items-center;
}

.navbar-pending {
  @apply w-screen h-[100vw] sm:w-[10.5rem] sm:h-[10.5rem] p-3 flex justify-center items-center;
}

.navbar-filter {
  @apply sm:relative;

  &__link {
    @apply text-white flex space-x-1 items-center font-semibold px-2 py-1 rounded-md hover:shadow-sm 
    shadow-blue-600/75;
  }

  &__img {
    @apply min-w-[1.5rem] w-6 object-contain inline-block drop-shadow-md;
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

  &:focus-within {
    .navbar-filter__link {
      @apply shadow-md;
    }

    .navbar-filter-items {
      @apply scale-y-100;
    }
  }

  &-items {
    @apply absolute flex flex-col bg-blue-400 sm:shadow-center sm:shadow-blue-600/75 sm:rounded-md transition-all overflow-y-auto overflow-x-hidden max-h-[86.5vh] sm:max-h-60 left-0 mt-2 sm:mt-0 sm:left-auto scale-y-0 origin-top;
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
      @apply bg-blue-600;
    }
  }

  &-item {
    &__img {
      @apply min-w-[1.25rem] w-5 object-contain inline-block drop-shadow-md rounded-md;
    }

    &__svg {
      @apply w-6;
    }

    &__link {
      @apply flex space-x-2 items-center font-semibold px-2 py-1 hover:bg-blue-600/50 text-white w-screen sm:w-44;

      &.active {
        @apply bg-blue-600;
      }

      &.sub-menu::after {
        // add arrow right icon
        @apply content-[''] absolute right-3 rotate-90 transition;
        border-top: 5px solid transparent;
        border-bottom: 5px solid transparent;
        border-left: 5px solid currentColor;
      }
    }
    &:focus-within {
      .sub-menu::after {
        @apply rotate-0;
      }
      .navbar-submenu-items {
        @apply flex flex-col;
      }
    }
  }
}

.navbar-submenu {
  &-items {
    @apply hidden;
  }

  &-item {
    @apply bg-blue-600/30 flex items-center w-full;

    &__link {
      @apply text-white flex space-x-2 items-center w-full font-semibold px-2 py-1 hover:bg-blue-600/50 pl-7 sm:pl-4;

      &.active {
        @apply bg-blue-600;
      }
    }

    &__img {
      @apply min-w-[1.25rem] w-5 object-contain inline-block rounded-sm;
    }

    &__svg {
      @apply w-6;
    }
  }
}

.nav-search {
  @apply inline-block mx-1 ml-3 flex-auto sm:flex-none relative;

  &__svg {
    @apply absolute mt-2 ml-2 text-blue-500;
  }

  &__input {
    @apply bg-blue-300 focus:bg-blue-200 disabled:bg-slate-500 py-1 px-2 rounded-lg transition-all hover:shadow-sm hover:shadow-blue-600/75 focus:shadow-md focus:shadow-blue-600/75 w-full text-gray-600 font-semibold placeholder:italic placeholder:text-blue-500 disabled:placeholder:text-blue-200 placeholder:font-normal pl-8 focus:outline-none;
  }
}
</style>
