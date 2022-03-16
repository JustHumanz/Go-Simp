<template>
  <div class="title">
    <span class="title__span">
      <font-awesome-icon icon="square-envelope" />
      Support
    </span>
  </div>
  <section>
    <h1 class="mt-2">About Vtbot</h1>
    <p>
      Vtbot (a.k.a Go-Simp) is a Open Source bot vtuber made by
      <a href="https://go.dev/" target="_blank" rel="noopener noreferrer"
        >Go Language</a
      >
      for sceduler, stats, share fanart, and more.
    </p>
  </section>
</template>

<script>
import { library } from "@fortawesome/fontawesome-svg-core"
import { faSquareEnvelope } from "@fortawesome/free-solid-svg-icons"

library.add(faSquareEnvelope)
export default {
  created() {
    document.title = "Support - Vtbot"
  },
}
</script>

<style lang="scss" scoped>
.title {
  @apply bg-blue-400 dark:bg-slate-500 py-3 w-full;

  &__span {
    @apply text-white text-2xl w-[90%] md:w-[70%] lg:w-[65%] font-semibold uppercase block mx-auto;
  }
  &__svg {
    @apply text-blue-200 dark:text-gray-200;
  }
}

section {
  @apply w-[90%] md:w-[70%] lg:w-[65%] mx-auto;
}

h1 {
  @apply text-2xl font-semibold;
}

a {
  @apply text-blue-500 dark:text-gray-300 hover:underline;
}
</style>
