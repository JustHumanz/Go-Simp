<script setup>
import { RouterLink, RouterView } from "vue-router"
import IconHome from "@/components/icons/IconHome.vue"
import "./index.css"
</script>

<template>
  <nav>
    <div class="navbar">
      <router-link to="/" class="icon"
        ><IconHome class="h-full" /> <span>Go-Simp</span></router-link
      >
      <div class="nav-menu">
        <a href="#" class="toggle-menu" onclick="return false"
          ><svg width="24" height="24" fill="none" aria-hidden="true">
            <path
              d="M12 6v.01M12 12v.01M12 18v.01M12 7a1 1 0 1 1 0-2 1 1 0 0 1 0 2Zm0 6a1 1 0 1 1 0-2 1 1 0 0 1 0 2Zm0 6a1 1 0 1 1 0-2 1 1 0 0 1 0 2Z"
              stroke="white"
              stroke-width="1.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            ></path></svg
        ></a>
        <ul>
          <li>
            <router-link
              to="/vtuber"
              :class="{ active: isActive }"
              class="nav-link"
              @click="resetFocus()"
              >Vtubers</router-link
            >
          </li>
          <li>
            <router-link to="/docs" class="nav-link" @click="resetFocus()"
              >Docs</router-link
            >
          </li>
          <li>
            <router-link to="/support" class="nav-link" @click="resetFocus()"
              >Support</router-link
            >
          </li>
          <li class="mobile-dashboard">
            <a
              href="https://web-admin.humanz.moe/login"
              target="_blank"
              class="nav-link"
              @click="resetFocus()"
              >Dashboard</a
            >
          </li>
        </ul>
      </div>

      <a
        class="dashboard-btn"
        href="https://web-admin.humanz.moe/login"
        target="_blank"
        >Dashboard</a
      >
    </div>
  </nav>

  <main class="mt-[4rem]">
    <RouterView />
  </main>
  <footer>
    <span>2022 &copy; JustHumanz</span>
    <div class="link">
      <!-- add link discord server -->
      <a href="https://discord.gg/ydWC5knbJT" target="_blank">Discord</a> |
      <!-- add link github -->
      <a href="https://github.com/JustHumanz/Go-Simp" target="_blank">Github</a>
    </div>
  </footer>
</template>

<script>
export default {
  data() {
    return {
      isActive: false,
    }
  },
  async created() {
    this.$watch(
      () => this.$route.params,
      async () => {
        this.isActive = this.$route.path.includes("/vtuber/")
      },

      { immediate: true }
    )
  },
  methods: {
    resetFocus() {
      this.$nextTick(() => {
        // get element .nav-link
        const navLinks = document.querySelectorAll(".nav-link")
        // looping nav-link
        navLinks.forEach((navLink) => {
          navLink.blur()
        })
      })
    },
  },
}
</script>

<style lang="scss">
nav {
  font-family: "Nunito", sans-serif;
  @apply h-16 bg-cyan-500 shadow-md shadow-cyan-500/50 fixed top-0 left-0 w-screen z-[11] flex justify-center;
  .navbar {
    @apply mx-4 h-full flex justify-between md:justify-around items-center w-full md:w-[90%] lg:w-[85%];

    .icon {
      @apply w-auto h-full p-3 cursor-pointer flex items-center;

      span {
        @apply ml-2 text-2xl font-semibold text-center uppercase sm:hidden text-white;
      }
    }

    .nav-menu {
      .toggle-menu {
        @apply block sm:hidden p-3;
      }

      &:focus-within ul {
        @apply visible;
      }

      ul {
        // make nav-menu items using tailwind
        @apply flex items-start sm:items-center flex-col sm:flex-row bg-cyan-500 rounded-md sm:bg-transparent w-[13rem] sm:w-auto absolute top-2 right-2 sm:top-auto sm:right-auto sm:relative invisible sm:visible shadow-center sm:shadow-none shadow-cyan-600/75;

        li {
          @apply w-full sm:w-auto sm:mx-1;

          &:first-child a {
            @apply rounded-t-md;
          }
          &:last-child a {
            @apply rounded-b-md;
          }

          a {
            @apply w-full inline-block cursor-pointer text-white py-2 sm:py-1 px-2 sm:rounded-md transition-shadow sm:hover:shadow-sm sm:hover:shadow-cyan-600/75 hover:bg-cyan-600/50 sm:hover:bg-transparent;

            &.router-link-active,
            &.active {
              @apply sm:shadow-md bg-cyan-600 sm:bg-transparent sm:shadow-cyan-600/75 font-semibold;
            }
          }

          &.mobile-dashboard {
            @apply block sm:hidden;

            a {
              @apply bg-rose-500 hover:bg-rose-600;
            }
          }
        }
      }
    }

    .dashboard-btn {
      @apply bg-rose-500 my-3 py-1 px-2 rounded-md font-semibold text-white cursor-pointer shadow-none transition-shadow hover:shadow-md hover:shadow-rose-500/75 hidden sm:block;
    }
  }
}
footer {
  @apply py-5 bg-cyan-500 font-semibold w-full flex justify-around;
  .link {
    @apply text-sm text-yellow-500;

    a {
      @apply text-black hover:text-orange-100 cursor-pointer;
    }
  }
}
</style>
