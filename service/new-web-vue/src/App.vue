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
  </main>
</template>

<style lang="scss">
nav {
  font-family: "Nunito", sans-serif;
  @apply h-16 bg-cyan-500 shadow-md shadow-cyan-500/50 fixed top-0 left-0 w-full z-10;
  .navbar {
    @apply mx-4 h-full flex justify-between items-center;

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
</style>
