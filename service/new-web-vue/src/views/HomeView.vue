<script setup></script>

<template>
  <section class="header-section">
    <h2>Welcome to Vtbot</h2>
    <p>Vtbot is a sceduler and sharing fan-art for vtubers.</p>
    <a href="/invite" target="_blank" class="btn-invite"> ADD ME! </a>
    <div class="text-orange-200">
      <a
        href="https://github.com/JustHumanz/Go-Simp"
        target="_blank"
        class="link-source"
      >
        Source Code
      </a>
      |
      <a
        href="https://discord.gg/ydWC5knbJT"
        target="_blank"
        class="link-source"
      >
        Join Server
      </a>
    </div>
  </section>
</template>

<style lang="scss">
.header-section {
  @apply flex flex-col items-center justify-center h-full py-52 bg-red-400;

  h2 {
    @apply text-7xl font-semibold text-center;
  }

  p {
    @apply text-xl text-center font-light text-gray-600;
  }

  .btn-invite {
    @apply bg-orange-400 my-3 py-2 px-4 text-xl rounded-md font-semibold text-white cursor-pointer;
  }
  .link-source {
    @apply text-sm text-center hover:text-orange-100 cursor-pointer;
  }
}
</style>
