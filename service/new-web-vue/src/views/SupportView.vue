<script setup>
import IconHome from "../components/icons/IconHome.vue"
</script>
<template>
  <div class="title">
    <span class="title__span">
      <font-awesome-icon icon="square-envelope" />
      Support
    </span>
  </div>
  <section>
    <div class="image-header mt-4">
      <IconHome />
    </div>
    <h1 class="mt-2">About Vtbot</h1>
    <p>
      Vtbot (a.k.a Go-Simp) is a Open Source bot vtuber made by
      <a href="https://go.dev/" target="_blank" rel="noopener noreferrer">
        <font-awesome-icon :icon="['fab', 'golang']" />
        Go Language</a
      >
      for ping about sceduler, stats, share fanart, and more in Discord Server.
      Over 400 vtubers added and more to come.
    </p>
    <h1 class="mt-2">Support This Bot</h1>
    <p>
      For now, we haven't received any form of payment yet, but you can vote
      <a
        href="https://top.gg/bot/721964514018590802"
        target="_blank"
        rel="noopener noreferrer"
        >here</a
      >
      for feedback. When you have a poblem or you want to contribute, you can
      request/join in our
      <a
        href="https://discord.gg/ydWC5knbJT"
        target="_blank"
        rel="noopener noreferrer"
      >
        <font-awesome-icon :icon="['fab', 'discord']" />
        Discord Server</a
      >
      or fork repository from
      <a
        href="https://github.com/JustHumanz/Go-Simp"
        target="_blank"
        rel="noopener noreferrer"
      >
        <font-awesome-icon :icon="['fab', 'github']" />
        Github</a
      >
    </p>
    <h1 class="mt-2">About Website</h1>
    <p class="mb-4">
      This is 3rd version of Vtbot website. This website using
      <a href="https://vuejs.org/" target="_blank" rel="noopener noreferrer">
        <font-awesome-icon :icon="['fab', 'vuejs']" class="fa-fw" />
        Vue.js</a
      >
      ,
      <a
        href="https://tailwindcss.com/"
        target="_blank"
        rel="noopener noreferrer"
        >Tailwind CSS</a
      >
      and some icon from
      <a
        href="https://fontawesome.com/"
        target="_blank"
        rel="noopener noreferrer"
      >
        <font-awesome-icon
          :icon="['fab', 'square-font-awesome']"
          class="fa-fw"
        />
        Font Awesome</a
      >
      and some smol holo gif and art by
      <a
        href="https://twitter.com/walfieee"
        target="_blank"
        rel="noopener noreferrer"
        ><font-awesome-icon :icon="['fab', 'twitter']" class="fa-fw" />
        @walfie</a
      >. This website create by
      <a
        href="https://github.com/yarndinasti"
        target="_blank"
        rel="noopener noreferrer"
        >Yarn</a
      >. You can support me
      <a href="https://yarn.my.id" target="_blank" rel="noopener noreferrer"
        >here</a
      >.
    </p>
  </section>
</template>

<script>
import { library } from "@fortawesome/fontawesome-svg-core"
import { faSquareEnvelope } from "@fortawesome/free-solid-svg-icons"
import {
  faVuejs,
  faSquareFontAwesome,
  faTwitter,
  faGolang,
} from "@fortawesome/free-brands-svg-icons"

library.add(faSquareEnvelope, faVuejs, faSquareFontAwesome, faTwitter, faGolang)
export default {
  created() {
    document.title = "Support - Vtbot"
  },
}
</script>

<style lang="scss" scoped>
.title {
  @apply w-full bg-blue-400 py-3 dark:bg-slate-500;

  &__span {
    @apply mx-auto block w-[90%] text-2xl font-semibold uppercase text-white md:w-[70%] lg:w-[65%];
  }
  &__svg {
    @apply text-blue-200 dark:text-gray-200;
  }
}

.image-header {
  @apply flex h-44 w-full justify-center py-3;

  :deep(svg) {
    @apply h-full object-contain;
  }
}

section {
  @apply mx-auto w-[90%] md:w-[70%] lg:w-[65%];
}

h1 {
  @apply text-2xl font-semibold;
}

a {
  @apply text-blue-500 hover:underline dark:text-gray-300;
}
</style>
