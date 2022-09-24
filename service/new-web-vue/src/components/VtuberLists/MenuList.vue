<template>
  <div class="vtuber-list-menu">
    <!-- Search -->
    <input
      type="text"
      class="vtuber-list-menu__search"
      ref="search_input"
      placeholder="Name Vtuber"
    />
    <!-- Refresh link -->
    <a href="#" onclick="return false" class="vtuber-list-menu__link">
      <font-awesome-icon
        class="fa-fw"
        icon="arrow-rotate-right"
        v-if="!refresh"
        @click="refresh = true"
      />
      <font-awesome-icon
        class="fa-fw animate-spin"
        icon="circle-notch"
        v-else
      />
    </a>
    <!-- Group/Agency -->
    <a href="#" onclick="return false" class="vtuber-list-menu__link">
      <!-- <font-awesome-icon class="fa-fw" icon="users" /> -->
      <img
        src="https://cdn.humanz.moe/holo.png"
        alt=""
        class="vtuber-list-menu__link-image"
      />
    </a>
    <!-- Filter -->
    <a href="#" onclick="return false" class="vtuber-list-menu__link">
      <font-awesome-icon class="fa-fw" icon="users" />
    </a>
    <!-- Sorting -->
    <a href="#" onclick="return false" class="vtuber-list-menu__link">
      <font-awesome-icon class="fa-fw" icon="users" />
    </a>
  </div>
</template>
<script>
import { library } from "@fortawesome/fontawesome-svg-core"
import {
  faMagnifyingGlass,
  faArrowRotateRight,
  faCircleNotch,
  faUsers,
  faUser,
} from "@fortawesome/free-solid-svg-icons"

library.add(
  faMagnifyingGlass,
  faArrowRotateRight,
  faCircleNotch,
  faUsers,
  faUser
)

import { useGroupStore } from "@/stores/groups"
import { useMemberStore } from "@/stores/members.js"

export default {
  data() {
    return {
      refresh: false,
    }
  },
}
</script>
<style lang="scss" scoped>
.vtuber-list-menu {
  @apply sticky mx-auto mb-3 mt-5 flex w-[95%] sm:w-[85%] md:w-[80%] lg:w-[75%];

  &__search {
    @apply rounded-full bg-rules-30 px-3 py-1 text-white transition-all duration-300 ease-in-out hover:brightness-110 focus:outline-none focus:brightness-125;
  }
  &__link {
    @apply ml-1 flex items-center justify-center rounded-full bg-rules-30 px-2 py-1 text-white transition-all duration-300 ease-in-out hover:brightness-110;

    &-image {
      @apply h-6 min-h-[1.5rem] object-contain;
    }
  }
}
</style>
