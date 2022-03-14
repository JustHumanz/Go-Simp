<script setup>
import ConfigView from "../components/DocsViews/Config.vue"
import TagingView from "../components/DocsViews/Taging.vue"
import UtilsView from "../components/DocsViews/Utils.vue"
</script>

<template>
  <div class="title">
    <span class="title__span">
      <font-awesome-icon
        :icon="['fas', 'circle-question']"
        class="title__svg"
      />
      Documentation
    </span>
  </div>
  <div class="tab" :class="{ fly: fly }">
    <a href="#" class="tab-link" onclick="return false"
      ><font-awesome-icon :icon="['fas', 'bars']" class="fa-fw tab-link__svg" />
      Menu</a
    >
    <ul class="tab-lists">
      <li class="tab-list">
        <router-link to="/docs/quick-start" class="tab-list__link"
          >Quick Start</router-link
        >
      </li>
      <li class="tab-list">
        <router-link to="/docs/configuration" class="tab-list__link"
          >Configuration</router-link
        >
      </li>
      <li class="tab-list">
        <router-link to="/docs/roles-and-taging" class="tab-list__link"
          >Roles and Taging</router-link
        >
      </li>
      <li class="tab-list">
        <router-link to="/docs/utilites" class="tab-list__link"
          >Utilites</router-link
        >
      </li>
    </ul>
  </div>
  <quick-start class="content" v-if="$route.params.page === 'quick-start'" />
</template>

<script>
import { library } from "@fortawesome/fontawesome-svg-core"

// import {

// } from "@fortawesome/free-brands-svg-icons"
import { faCircleQuestion, faBars } from "@fortawesome/free-solid-svg-icons"
import QuickStart from '../components/DocsViews/QuickStart.vue'

library.add(faCircleQuestion, faBars)

// change title
export default {
  components: { QuickStart },
  data() {
    return {
      docs: null,
      fly: false,
    }
  },
  mounted() {
    document.title = "Documentation - Vtbot"
  },
  async created() {

    window.addEventListener("scroll", () => {
      let bottomOfWindow = Math.ceil(window.scrollY) >= 55

      this.fly = bottomOfWindow
    })
  },
  computed: {
    isActive() {
      return this.$route.path === "/docs"
    },
  },
  methods: {
    resetFocus() {
      this.$refs.search.focus()
    },
    changeTab() {
      const hash = this.$route.hash.replace("#", "")

      console.log(hash)
      switch (hash) {
        case "config":
          this.tab = 1
          break
        case "tags":
          this.tab = 2
          break
        case "utils":
          this.tab = 3
          break
        case "data":
          this.tab = 4
          break
      }
    },
  },
}
</script>

<style lang="scss" scoped>
.title {
  @apply text-2xl font-semibold uppercase bg-blue-400 dark:bg-slate-500 py-3 w-full flex flex-wrap relative;

  &__span {
    @apply w-[90%] md:w-[70%] lg:w-[65%] mx-auto text-white;
  }
  &__svg {
    @apply text-blue-200 dark:text-gray-200;
  }
}

.fly {
  position: fixed !important;
  @apply top-16;

  .tab-link {
    @apply shadow-md bg-blue-400 text-white dark:bg-slate-500 dark:text-white;
  }

  .tab-lists {
    @apply shadow-md md:shadow-none z-[2];
  }
}

.tab {
  @apply flex-col md:py-2 md:border-r-2 border-blue-300 dark:border-white w-full md:w-48 md:mx-0 md:ml-3 md:absolute md:h-screen;

  &-link {
    @apply text-xl py-3 w-full inline-block md:hidden transition-all font-semibold px-6 bg-blue-300 dark:bg-slate-400 dark:text-gray-700;

    &__svg {
      @apply mr-2;
    }
  }

  &:focus-within {
    .tab-lists {
      @apply scale-y-100;
    }
  }

  &-lists {
    @apply scale-y-0 origin-top md:scale-y-100 flex flex-col w-full absolute md:static md:z-auto bg-blue-100 dark:bg-slate-500 md:bg-transparent dark:md:bg-transparent transition-all duration-200 ease-in-out;
  }

  &-list {
    @apply text-blue-400 dark:text-slate-200;

    &__link {
      @apply w-full h-full inline-block px-3 py-1 md:rounded-l-full hover:bg-blue-200 dark:hover:bg-slate-700;
      transition: all 0.2s ease-in-out;

      &.router-link-active {
        @apply bg-blue-300 dark:text-slate-800 dark:bg-white text-white font-semibold;
      }
    }
  }
}

.content {
  @apply px-3 py-4 md:pl-[13.5rem];
}
</style>
