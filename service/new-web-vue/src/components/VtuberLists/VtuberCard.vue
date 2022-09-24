<template>
  <div class="vtuber-cards m-5 w-44">
    <img src="/assets/smolame.jpg" alt="" class="vtuber-cards-cover__images" />
    Name
  </div>
</template>
<script></script>
<style lang="scss" scoped>
.vtuber-cards {
  @apply overflow-hidden rounded-md bg-red-500 transition-all duration-300 ease-in-out hover:rounded-lg;

  &-cover {
    &__images {
      @apply aspect-square w-full object-contain;
    }
  }
}
</style>
