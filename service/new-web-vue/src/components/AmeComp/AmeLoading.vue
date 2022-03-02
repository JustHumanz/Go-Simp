<template>
  <section class="ame-loading">
    <img
      :src="`/src/assets/loading/${Math.floor(Math.random() * 7)}.gif`"
      class="ame-loading__img"
    />
  </section>
</template>

<style lang="scss" scoped>
.ame-loading {
  @apply flex flex-col items-center justify-center w-full h-40;

  &__img {
    @apply h-32 object-contain;
  }
}
</style>
