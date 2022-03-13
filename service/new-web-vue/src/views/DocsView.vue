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
      ><font-awesome-icon :icon="['fas', 'bars']" class="fa-fw tab-link__svg" /> Menu</a
    >
    <ul class="tab-lists">
      <li class="tab-list">
        <a href="#" class="tab-list__link">Quick Setup</a>
      </li>
      <li class="tab-list">
        <a href="#" class="tab-list__link">Configuration</a>
      </li>
      <li class="tab-list">
        <a href="#" class="tab-list__link">Roles and Taging</a>
      </li>
      <li class="tab-list">
        <a href="#" class="tab-list__link">Utilites</a>
      </li>
    </ul>
  </div>
  <section class="content">
    
  </section>
</template>

<script>
import { library } from "@fortawesome/fontawesome-svg-core"

// import {

// } from "@fortawesome/free-brands-svg-icons"
import { faCircleQuestion, faBars } from "@fortawesome/free-solid-svg-icons"

library.add(faCircleQuestion, faBars)

// change title
export default {
  data() {
    return {
      tab: 1,
      fly: false,
    }
  },
  mounted() {
    document.title = "Documentation - Vtbot"
  },
  created() {
    this.changeTab()
    console.log(FileSystem.root)
    window.addEventListener("scroll", () => {
      let bottomOfWindow = Math.ceil(window.scrollY) >= 55

      this.fly = bottomOfWindow
      // if (bottomOfWindow) {
      //   console.log("scroll")
      // } else {
      //   console.log("no scroll")
      // }
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
  @apply text-2xl font-semibold uppercase bg-blue-400 py-3 w-full flex flex-wrap relative;

  &__span {
    @apply w-[90%] md:w-[70%] lg:w-[65%] mx-auto text-white;
  }
  &__svg {
    @apply text-blue-200;
  }
}

.fly {
    position: fixed !important;
    @apply top-16;

    .tab-link{
      @apply shadow-md bg-blue-400 text-white;
    }

    .tab-lists {
      @apply shadow-md md:shadow-none z-[2];
    }
}

.message {
  @apply bg-yellow-400 text-white px-4 py-2 mx-2 my-1 rounded-md font-semibold;
}
.tab {
  @apply flex-col md:py-2 md:border-r-2 border-blue-300 w-full md:w-48 md:mx-0 md:ml-3 md:absolute md:h-screen;

  &-link {
    @apply text-xl py-3 w-full inline-block md:hidden transition-all font-semibold px-6 bg-blue-300;

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
    @apply scale-y-0 origin-top md:scale-y-100 flex flex-col w-full absolute md:static md:z-auto bg-blue-100 md:bg-transparent transition-all duration-200 ease-in-out;
  }

  &-list {
    @apply px-3 py-1 text-blue-400 hover:bg-blue-100 md:rounded-l-full;
    transition: all 0.2s ease-in-out;

    &.active {
      @apply bg-blue-300 text-white font-semibold;
    }

    &__link {
      @apply w-full h-full inline-block;
    }
  }
}

.content {
  @apply px-3 md:pl-[13.5rem];
}
</style>
