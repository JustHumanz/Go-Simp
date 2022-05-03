<template>
  <section class="ame-loading">
    <img
      draggable="false"
      :src="`/assets/loading/${Math.floor(Math.random() * 7)}.gif`"
      class="ame-loading__img"
    />
  </section>
</template>

<style lang="scss" scoped>
.ame-loading {
  @apply flex h-40 w-full flex-col items-center justify-center;

  &__img {
    @apply h-32 object-contain;
  }
}
</style>
