<script setup>
import { RouterLink, RouterView } from "vue-router"
import IconHome from "@/components/icons/IconHome.vue"
import "./index.css"

import { useGroupStore } from "@/stores/groups"
import { onMounted } from "vue"

onMounted(() => {
  useGroupStore().fetchGroups()
})
</script>

<template>
  <nav class="nav">
    <div class="navbar">
      <router-link to="/" class="navbar-icon"
        ><IconHome class="h-full" />
        <span class="navbar-icon__span"> Go-Simp </span></router-link
      >
      <div class="navbar-group">
        <a
          href="#"
          class="navbar-toggle navbar-toggle-dark"
          @click="darkMode"
          onclick="return false"
        >
          <font-awesome-icon
            icon="moon"
            class="navbar-toggle-dark__icon fa-fw"
            v-if="theme === 'dark'"
          />
          <font-awesome-icon
            icon="sun"
            class="navbar-toggle-dark__icon fa-fw"
            v-else-if="theme === 'light'"
          />
          <font-awesome-icon
            icon="circle-half-stroke"
            class="navbar-toggle-dark__icon fa-fw"
            v-else
          />
        </a>

        <div class="navbar-menu">
          <a href="#" class="navbar-toggle" onclick="return false">
            <font-awesome-icon
              icon="ellipsis-vertical"
              class="navbar-toggle__icon fa-fw"
            />
          </a>
          <ul class="navbar-items">
            <li class="navbar-item">
              <router-link
                to="/vtubers"
                :class="{ active: isActiveVtuber }"
                class="navbar-link"
                >Vtubers</router-link
              >
            </li>
            <li class="navbar-item">
              <router-link
                to="/docs"
                :class="{ active: isActiveDocs }"
                class="navbar-link"
                >Docs</router-link
              >
            </li>
            <li class="navbar-item">
              <router-link to="/support" class="navbar-link"
                >Support</router-link
              >
            </li>
            <li class="navbar-item">
              <a
                href="https://top.gg/bot/721964514018590802"
                target="_blank"
                class="navbar-link"
                >Vote</a
              >
            </li>
            <li class="navbar-item mobile-menu invite">
              <a
                href="https://discord.com/oauth2/authorize?client_id=721964514018590802&permissions=456720&scope=bot%20applications.commands"
                target="_blank"
                class="navbar-link"
                ><font-awesome-icon
                  icon="right-to-bracket"
                  class="fa-fw mr-1"
                />
                Invite Me</a
              >
            </li>
            <li class="navbar-item mobile-menu dashboard">
              <a
                href="https://web-admin.humanz.moe/login"
                target="_blank"
                class="navbar-link"
                ><font-awesome-icon icon="gauge-simple" class="fa-fw mr-1" />
                Dashboard</a
              >
            </li>
          </ul>
        </div>
      </div>

      <div class="navbar-buttons">
        <a
          class="navbar-buttons__button invite group"
          href="https://discord.com/oauth2/authorize?client_id=721964514018590802&permissions=456720&scope=bot%20applications.commands"
          target="_blank"
          ><font-awesome-icon icon="right-to-bracket" class="fa-fw" />
          <span
            class="navbar-buttons__hover group-hover:!scale-100 group-hover:!opacity-100"
            >Invite Me</span
          >
        </a>
        <a
          class="navbar-buttons__button dashboard group"
          href="https://web-admin.humanz.moe/login"
          target="_blank"
          ><font-awesome-icon icon="gauge-simple" class="fa-fw" />
          <span
            class="navbar-buttons__hover group-hover:!scale-100 group-hover:!opacity-100"
            >Dashboard</span
          >
        </a>
        <!-- Add Dark mode button -->
        <a
          href="#"
          class="navbar-buttons__button dark-mode group"
          @click="darkMode"
          onclick="return false"
        >
          <font-awesome-icon
            icon="moon"
            class="dark-mode-btn__svg fa-fw"
            v-if="theme === 'dark'"
          />
          <font-awesome-icon
            icon="sun"
            class="dark-mode-btn__svg fa-fw"
            v-else-if="theme === 'light'"
          />
          <font-awesome-icon
            icon="circle-half-stroke"
            class="dark-mode-btn__svg fa-fw"
            v-else
          />
          <span
            class="navbar-buttons__hover group-hover:!scale-100 group-hover:!opacity-100"
          >
            {{
              theme === "dark"
                ? "Dark Mode"
                : theme === "light"
                ? "Light Mode"
                : "Auto"
            }}
          </span>
        </a>
      </div>
    </div>
  </nav>

  <main id="main">
    <RouterView />
  </main>
  <footer class="footer relative z-[1]">
    <div class="footer-page">
      <h2 class="footer-page__title">Vtbot</h2>
      <h4 class="footer-page__subtitle">Bot for Vtubers</h4>
    </div>
    <div class="footer-links">
      <h4 class="footer-links__title">Information</h4>
      <ul class="footer-lists">
        <li class="footer-list">
          <router-link to="/docs" class="footer-list__link">Docs</router-link>
        </li>
        <li class="footer-list">
          <router-link to="/support" class="footer-list__link"
            >Support</router-link
          >
        </li>
        <!-- Github link -->
        <li class="footer-list">
          <a
            href="https://github.com/JustHumanz/Go-Simp"
            class="footer-list__link"
            target="_blank"
          >
            <font-awesome-icon
              :icon="['fab', 'github']"
              class="footer-list__icon fa-fw"
            ></font-awesome-icon>
            Github
          </a>
        </li>
        <!-- Discord link -->
        <li class="footer-list">
          <a
            href="https://discord.gg/ydWC5knbJT"
            target="_blank"
            class="footer-list__link"
          >
            <font-awesome-icon
              :icon="['fab', 'discord']"
              class="footer-list__icon fa-fw"
            ></font-awesome-icon>
            Discord Server
          </a>
        </li>
      </ul>
    </div>
    <div class="footer-credit">
      This Website designed and developed by
      <a
        href="https://github.com/yarndinasti"
        target="_blank"
        class="footer-credit__link"
        >Yarn</a
      >
    </div>
  </footer>
</template>

<script>
// Add Font Awesome for Discord and Github icons
import { library } from "@fortawesome/fontawesome-svg-core"
import { faDiscord, faGithub } from "@fortawesome/free-brands-svg-icons"
import {
  faMoon,
  faSun,
  faCircleHalfStroke,
  faGaugeSimple,
  faEllipsisVertical,
  faRightToBracket,
} from "@fortawesome/free-solid-svg-icons"
import { watch } from "@vue/runtime-core"
library.add(
  faDiscord,
  faGithub,
  faMoon,
  faSun,
  faGaugeSimple,
  faCircleHalfStroke,
  faEllipsisVertical,
  faRightToBracket
)

export default {
  data() {
    return {
      isActiveVtuber: false,
      isActiveDocs: false,
      theme: null,
    }
  },
  async mounted() {
    this.getClickMenu()
    this.changeView()

    // get viewport height
    window.addEventListener("resize", () => this.changeView())

    this.theme = localStorage.getItem("theme")

    this.$watch(
      () => this.$route.params,
      async () => {
        this.isActiveVtuber = this.$route.path.includes("/vtuber")
        this.isActiveDocs = this.$route.path.includes("/docs")
      },

      { immediate: true }
    )

    this.$watch(
      () => this.theme,
      async () => {
        switch (this.theme) {
          case "dark":
            document.documentElement.classList.add("dark")
            break
          case "light":
            document.documentElement.classList.remove("dark")
            break
          default:
            this.autoDark()
        }
      },
      { immediate: true }
    )
  },
  methods: {
    getClickMenu() {
      document.body.addEventListener("click", (e) => {
        if (e.target.id === "router-link") {
          e.preventDefault()
          this.$router.push(e.target.getAttribute("href"))
        }

        if (e.target.closest(".navbar-link")) {
          e.target.closest(".navbar-link").blur()
        }
      })
    },

    changeView() {
      // get viewport height
      const main_height = window.innerHeight - 204
      const mainEl = document.getElementById("main")

      // set min-height for main element
      if (mainEl) mainEl.style.minHeight = `${main_height}px`
    },

    darkMode() {
      switch (this.theme) {
        case "dark":
          localStorage.setItem("theme", "light")
          this.theme = "light"
          break
        case "light":
          localStorage.removeItem("theme")
          this.theme = null
          break
        default:
          localStorage.setItem("theme", "dark")
          this.theme = "dark"
          break
      }
      console.log(this.theme)
    },
    autoDark() {
      if (this.theme) return

      const darkSheme = window.matchMedia("(prefers-color-scheme: dark)")

      if (darkSheme.matches) {
        document.documentElement.classList.add("dark")
      } else {
        document.documentElement.classList.remove("dark")
      }

      darkSheme.onchange = async () => {
        if (darkSheme.matches) {
          document.documentElement.classList.add("dark")
        } else {
          document.documentElement.classList.remove("dark")
        }
      }
    },
  },
}
</script>

<style lang="scss" scoped>
.nav {
  font-family: "Nunito", sans-serif;
  @apply sticky top-0 left-0 z-[20] flex h-16 w-screen select-none justify-center bg-cyan-500 shadow-md shadow-cyan-500/50 dark:bg-slate-700 dark:shadow-slate-700/50;
  .navbar {
    @apply mx-4 flex h-full w-full items-center justify-between md:w-[90%] md:justify-around lg:w-[85%];

    &-group {
      @apply flex;
    }

    &-menu {
      @apply relative;
    }

    &-icon {
      @apply flex h-full w-auto cursor-pointer items-center p-3;

      &__span {
        @apply ml-2 text-center text-2xl font-semibold uppercase text-white sm:hidden;
      }
    }
    .navbar-toggle {
      @apply block p-3 sm:hidden;

      &__icon,
      &-dark__icon {
        @apply text-white;
      }
    }

    &-menu:focus-within .navbar-items {
      @apply scale-100;
    }

    .navbar-items {
      // make nav-menu items using tailwind
      @apply absolute top-3 right-3 flex w-[13rem] origin-top-right scale-0 flex-col items-start overflow-hidden rounded-md bg-cyan-500 shadow-md shadow-black/50 transition duration-200 ease-in-out dark:bg-slate-600 sm:static sm:top-auto sm:right-auto sm:w-auto sm:scale-100 sm:flex-row sm:items-center sm:overflow-visible sm:bg-transparent sm:shadow-none sm:transition-none sm:dark:bg-slate-700;
    }

    .navbar-item {
      @apply w-full sm:mx-1 sm:w-auto;

      .navbar-link {
        @apply inline-block w-full cursor-pointer py-2 px-2 font-semibold text-white transition hover:bg-black/20 sm:rounded-md sm:py-1 hover:sm:-translate-y-px sm:hover:bg-transparent sm:hover:shadow-sm sm:hover:shadow-cyan-600/75 sm:hover:dark:bg-transparent sm:dark:hover:shadow-slate-200/75;

        &.router-link-active,
        &.active {
          @apply bg-cyan-600 dark:bg-slate-800 sm:-translate-y-0.5 sm:bg-transparent sm:shadow-md sm:shadow-cyan-600/75 dark:sm:bg-transparent sm:dark:shadow-slate-200/75;
        }
      }

      &.mobile-menu {
        @apply block sm:hidden;

        &.dashboard {
          @apply bg-rose-500 dark:bg-slate-800/50;
        }
        &.invite {
          @apply bg-indigo-500 dark:bg-slate-800/20;
        }
      }
    }

    &-buttons {
      @apply hidden sm:block;

      &__hover {
        @apply absolute top-9 left-1/2 inline-block w-max -translate-x-1/2 scale-0 rounded-md px-2 py-1 opacity-0 shadow-md transition;
      }

      &__button {
        @apply relative my-3 cursor-pointer rounded-md py-1 px-2 font-semibold text-white shadow-none transition-shadow hover:-translate-y-0.5 hover:shadow-md;

        &.dashboard {
          @apply mr-2 bg-rose-500 hover:shadow-rose-500/75 dark:bg-slate-900 hover:dark:shadow-slate-100/75;

          .navbar-buttons__hover {
            @apply bg-rose-500 shadow-rose-500/75 dark:bg-slate-900 dark:shadow-slate-100/75;
          }
        }
        &.invite {
          @apply mr-2 bg-indigo-500 hover:shadow-indigo-500/75 dark:bg-slate-800 hover:dark:shadow-slate-100/75;

          .navbar-buttons__hover {
            @apply bg-indigo-500 shadow-indigo-500/75 dark:bg-slate-800 dark:shadow-slate-100/75;
          }
        }

        &.dark-mode {
          @apply bg-orange-500 hover:shadow-sky-500/75 dark:bg-slate-800 hover:dark:shadow-slate-100/75;

          .navbar-buttons__hover {
            @apply bg-orange-500 shadow-sky-500/75 dark:bg-slate-800 dark:shadow-slate-100/75;
          }
        }
      }
    }
  }
}
.footer {
  @apply grid w-full justify-around bg-cyan-500 pt-3 font-semibold dark:bg-slate-700;

  grid-template-columns: 1fr 1fr;
  grid-template-areas:
    "page links"
    "credit credit";

  &-page {
    @apply flex flex-col justify-center pl-5 sm:px-4 md:pl-8 lg:pl-12 xl:pl-20;
    grid-area: page;

    &__title {
      @apply text-3xl font-semibold text-cyan-200 dark:text-slate-300;
    }
    &__subtitle {
      @apply text-white;
    }
  }

  &-links {
    @apply pl-5 md:pl-8 lg:pl-12 xl:pl-20;
    grid-area: links;
    &__title {
      @apply text-xl font-semibold text-cyan-200 dark:text-slate-300;
    }
  }

  &-lists {
    @apply grid;

    @media (min-width: 640px) {
      @apply grid-flow-col;
      grid-template-rows: repeat(2, 1fr);
    }
  }

  &-list {
    @apply w-full;

    &__link {
      @apply inline-block w-full text-white hover:text-cyan-100 dark:hover:text-slate-300;
    }
  }

  .footer-credit {
    grid-area: credit;
    @apply mt-3 bg-cyan-700 py-2 text-center text-cyan-50 dark:bg-slate-800;
    &__link {
      @apply text-white hover:text-cyan-100 dark:hover:text-slate-300;
    }
  }
}
</style>
