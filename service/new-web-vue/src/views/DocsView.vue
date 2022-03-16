<script setup>
import DocsRender from "../components/DocsRender.vue"
import NotFound from "./NotFound.vue"
</script>

<template>
<div v-if="markdowns && markdowns[`./../components/docs/${$route.params.page}.md`]">
  <div class="header-title">
    <div class="title">
      <span class="title__span">
        <font-awesome-icon
          :icon="['fas', 'circle-question']"
          class="title__svg"
        />
        Documentation
      </span>
      <a
        :href="`https://github.com/JustHumanz/Go-Simp/blob/new-web-lets-go/service/new-web-vue/src/components/docs/${$route.params.page}.md`"
        target="_blank"
        rel="noopener noreferrer"
        class="edit-github"
      >
        <font-awesome-icon :icon="['fab', 'github']" class="edit-github__svg" />
        Edit on GitHub</a
      >
    </div>
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
        <router-link to="/docs/utilities" class="tab-list__link"
          >Utilities</router-link
        >
      </li>
    </ul>
  </div>
  <docs-render class="content" :page="$route.params.page" />
  </div>
  <div v-else>
    <not-found />
  </div>
</template>

<script>
const mdfiles = import.meta.glob("./../components/docs/*.md", { assert: { type: "raw" } })
import { library } from "@fortawesome/fontawesome-svg-core"

import { faGithub } from "@fortawesome/free-brands-svg-icons"
import { faCircleQuestion, faBars } from "@fortawesome/free-solid-svg-icons"

library.add(faCircleQuestion, faBars, faGithub)

// change title
export default {
  data() {
    return {
      markdowns: null,
      fly: false,
    }
  },
  mounted() {
    this.markdowns = mdfiles

    this.$watch(
      () => this.$route.params.page,
      () => {
        if (!this.$route.params.page) return
        document.title = `${this.$route.params.page
          .split("-")
          .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
          .join(" ")} - Documentation`
      },
      {
        immediate: true,
      }
    )
  },
  async created() {
    window.addEventListener("scroll", () => {
      let bottomOfWindow = Math.ceil(window.scrollY) >= 55

      this.fly = bottomOfWindow
    })
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
.header-title {
  @apply bg-blue-400 dark:bg-slate-500 py-3 w-full;
}

.title {
  @apply w-[90%] md:w-[70%] lg:w-[65%] mx-auto flex flex-wrap justify-between;

  &__span {
    @apply text-white text-2xl font-semibold uppercase;
  }
  &__svg {
    @apply text-blue-200 dark:text-gray-200;
  }
}

.edit-github {
  @apply px-2 py-1 inline-block rounded-md shadow-sm -translate-y-px hover:shadow-md hover:-translate-y-0.5 shadow-blue-600/75 dark:shadow-slate-300/50 text-sm font-semibold text-white transition duration-200 ease-in-out;
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
