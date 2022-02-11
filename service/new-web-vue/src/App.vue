<script setup>
import { RouterLink, RouterView } from "vue-router"
import IconHome from "@/components/IconHome.vue"
import "./index.css"
</script>

<template>
  <nav>
    <div class="navbar">
      <a to="/" class="icon"><IconHome class="h-full" /></a>
      <ul class="nav-menu">
        <li><a to="">Vtubers</a></li>
        <li><a to="">Docs</a></li>
        <li><a to="">Support</a></li>
      </ul>
      <a class="dashboard-btn">Dashboard</a>
    </div>
  </nav>

  <main class="mt-[4rem]">
    <RouterView />
    <footer>
      <span>2022 &copy; JustHumanz</span>
      <div class="link">
        <!-- add link discord server -->
        <a href="https://discord.gg/ydWC5knbJT" target="_blank">Discord</a> |
        <!-- add link github -->
        <a href="https://github.com/JustHumanz/Go-Simp" target="_blank"
          >Github</a
        >
      </div>
    </footer>
  </main>
</template>

<style lang="scss">
nav {
  font-family: "Nunito", sans-serif;
  @apply h-16 bg-cyan-500 shadow-md shadow-cyan-500/50 fixed top-0 left-0 w-full z-10 flex justify-center;
  .navbar {
    @apply mx-4 h-full flex justify-between md:justify-around items-center w-full md:w-[90%] lg:w-[85%];

    .icon {
      @apply w-auto h-full p-3 cursor-pointer;
    }

    .nav-menu {
      // make nav-menu items using tailwind
      @apply flex items-center;

      li {
        @apply p-3 text-white hover:text-cyan-100 cursor-pointer;
      }
    }

    .dashboard-btn {
      @apply bg-orange-400 my-3 py-1 px-2 rounded-md font-semibold text-white cursor-pointer hover:shadow-md shadow-orange-400/75;
    }
  }
}
footer {
  @apply py-5 bg-cyan-500 shadow-md shadow-cyan-500 bottom-0 left-0 w-full flex justify-around;
  font-family: "Nunito", sans-serif;
  .link {
    a {
      @apply text-sm hover:text-orange-100 cursor-pointer;
    }
  }
}
</style>
