<script setup>
import DocsRender from "../components/DocsRender.vue"
import NotFound from "./NotFound.vue"
</script>

<template>
  <div
    v-if="markdowns && markdowns.find((md) => md.slug === $route.params.page)"
  >
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
          :href="`https://github.com/JustHumanz/Go-Simp/blob/master/service/new-web-vue/src/components/docs/${$route.params.page}.md`"
          target="_blank"
          rel="noopener noreferrer"
          class="edit-github"
        >
          <font-awesome-icon
            :icon="['fab', 'github']"
            class="edit-github__svg"
          />
          Edit on GitHub</a
        >
      </div>
    </div>
    <div class="tab" :class="{ fly: fly }">
      <a href="#" class="tab-link" onclick="return false"
        ><font-awesome-icon
          :icon="['fas', 'bars']"
          class="fa-fw tab-link__svg"
        />
        Menu</a
      >
      <ul class="tab-lists">
        <li class="tab-list">
          <router-link
            v-for="md in markdowns"
            :key="md.slug"
            :to="`/docs/${md.slug}`"
            class="tab-list__link"
            >{{ md.title }}</router-link
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
const mdfiles = import.meta.glob("./../components/docs/*.md", {
  as: "raw",
})
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
      title: "",
      menuDocs: false,
    }
  },
  mounted() {
    this.markdowns = Object.keys(mdfiles).map((path) => {
      const slug = path.replace("./../components/docs/", "").replace(".md", "")
      return {
        slug,
        title: this.convertToTitle(slug),
      }
    })
    // sort quick-start slug to the top
    this.markdowns.sort((a, b) => {
      if (a.slug === "quick-start") {
        return -1
      }
      if (b.slug === "quick-start") {
        return 1
      }
      return 0
    })

    this.$watch(
      () => this.$route.params.page,
      () => {
        if (!this.$route.params.page) return
        this.title = this.convertToTitle(this.$route.params.page)

        document.title = `${this.title} - Documentation`
      },
      {
        immediate: true,
      }
    )

    document.body.addEventListener("click", (e) => {
      if (e.target.closest(".tab-link")) {
        const docsMenu = e.target.closest(".tab-link")

        this.menuDocs = !this.menuDocs
        if (!this.menuDocs && document.activeElement === docsMenu) {
          docsMenu.blur()
        }
      } else if (e.target.closest(".tab-list__link")) {
        document.activeElement.blur()
        this.menuDocs = false
      } else this.menuDocs = false
    })

    document.body.addEventListener("blur", () => {
      if (!this.menuDocs) return
      document.activeElement.blur()
      this.menuDocs = false
    })
  },
  async created() {
    window.addEventListener("scroll", () => {
      let bottomOfWindow = Math.ceil(window.scrollY) >= 55

      this.fly = bottomOfWindow
    })
  },
  methods: {
    convertToTitle(str) {
      return str
        .split("-")
        .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
        .join(" ")
    },
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
  @apply w-full bg-blue-400 py-3 dark:bg-slate-500;
}

.title {
  @apply mx-auto flex w-[90%] flex-wrap justify-between md:w-[70%] lg:w-[65%];

  &__span {
    @apply text-2xl font-semibold uppercase text-white;
  }
  &__svg {
    @apply text-blue-200 dark:text-gray-200;
  }
}

.edit-github {
  @apply inline-block -translate-y-px rounded-md px-2 py-1 text-sm font-semibold text-white shadow-sm shadow-blue-600/75 transition duration-200 ease-in-out hover:-translate-y-0.5 hover:shadow-md dark:shadow-slate-300/50;
}

.fly {
  position: fixed !important;
  @apply top-16;

  .tab-link {
    @apply bg-blue-400 text-white shadow-md dark:bg-slate-500 dark:text-white;
  }

  .tab-lists {
    @apply z-[2] shadow-md md:shadow-none;
  }
}

.tab {
  @apply w-full flex-col border-blue-300 dark:border-white md:absolute md:mx-0 md:ml-3 md:h-screen md:w-48 md:border-r-2 md:py-2;

  &-link {
    @apply inline-block w-full bg-blue-300 py-3 px-6 text-xl font-semibold transition-all dark:bg-slate-400 dark:text-gray-700 md:hidden;

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
    @apply absolute flex w-full origin-top scale-y-0 flex-col bg-blue-100 transition-all duration-200 ease-in-out dark:bg-slate-500 md:static md:z-auto md:scale-y-100 md:bg-transparent dark:md:bg-transparent;
  }

  &-list {
    @apply text-blue-400 dark:text-slate-200;

    &__link {
      @apply inline-block h-full w-full px-3 py-1 hover:bg-blue-200 dark:hover:bg-slate-700 md:rounded-l-full;
      transition: all 0.2s ease-in-out;

      &.router-link-active {
        @apply bg-blue-300 font-semibold text-white dark:bg-white dark:text-slate-800;
      }
    }
  }
}

.content {
  @apply px-3 py-4 md:pl-[13.5rem];
}
</style>
