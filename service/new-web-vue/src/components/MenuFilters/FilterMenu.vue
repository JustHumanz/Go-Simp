<script setup>
import GroupsMenu from "./GroupsMenu.vue"
</script>

<template>
  <nav class="list-nav">
    <ul class="navbar-filters">
      <li class="navbar-filter">
        <GroupsMenu :groupid="$route.params?.id" />
      </li>
      <li class="navbar-filter">
        <a href="#" class="navbar-filter__link" onclick="return false"
          >Filters</a
        >
      </li>
      <li class="navbar-filter">
        <a href="#" class="navbar-filter__link" onclick="return false">Sort</a>
      </li>
    </ul>
    <div class="nav-search">
      <font-awesome-icon
        icon="magnifying-glass"
        class="fa-fw fa-md nav-search__svg"
      />
      <input type="text" class="nav-search__input" />
    </div>
  </nav>
</template>

<script>
import { library } from "@fortawesome/fontawesome-svg-core"
import { faMagnifyingGlass } from "@fortawesome/free-solid-svg-icons"

library.add(faMagnifyingGlass)
</script>

<style lang="scss">
.list-nav {
  @apply bg-blue-400 fixed top-16 py-2 px-5 w-screen flex flex-wrap-reverse items-center sm:justify-around justify-center z-10;
}

.navbar-filters {
  @apply flex space-x-1 sm:space-x-2 first:mt-2 xs:first:mt-0 items-center;
}

.navbar-filter {
  @apply relative;

  &__link {
    @apply text-white flex space-x-1 items-center font-semibold px-2 py-1 rounded-md hover:shadow-sm hover:shadow-blue-600/75;
  }

  &__img {
    @apply min-w-[1.5rem] max-w-[1.5rem] object-contain inline-block drop-shadow-md;
  }

  &-items {
    @apply flex flex-col absolute;
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