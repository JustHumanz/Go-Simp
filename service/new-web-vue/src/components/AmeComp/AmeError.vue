<template>
  <section class="ame-error">
    <img :src="`/assets/smolame/${img}.png`" alt="" class="ame-error-image" />
    <div class="ame-error-text">
      <h2 class="ame-error-text__title">
        <font-awesome-icon
          icon="circle-exclamation"
          class="fa-fw ame-error-text__icon ml-2 text-red-600"
          v-if="type === `error`"
        />
        <font-awesome-icon
          icon="exclamation-triangle"
          class="fa-fw ame-error-text__icon ml-2 text-yellow-400"
          v-else-if="type === `warning`"
        />
        <span>{{ title }}</span>
      </h2>
      <p v-html="description" class="ame-error-text__description"></p>
      <small class="ame-error-text__err-msg" v-if="message">
        {{ message }}
      </small>
    </div>
  </section>
</template>

<script>
import { library } from "@fortawesome/fontawesome-svg-core"
import {
  faTriangleExclamation,
  faCircleExclamation,
} from "@fortawesome/free-solid-svg-icons"

library.add(faTriangleExclamation, faCircleExclamation)

export default {
  props: {
    img: {
      type: String,
    },
    type: {
      type: String,
    },
    title: {
      type: String,
    },
    description: {
      type: String,
    },
    message: {
      type: String,
    },
  },
}
</script>

<style lang="scss">
.ame-error-text__link {
  @apply text-gray-700 underline decoration-dotted underline-offset-1 hover:text-blue-400 hover:decoration-2 dark:text-slate-200 dark:hover:text-slate-400;
}
</style>

<style lang="scss" scoped>
.ame-error {
  @apply flex h-screen w-full flex-col items-center justify-center md:flex-row;

  &-image {
    @apply h-24 w-24 object-contain;
  }

  &-text {
    @apply max-w-[29.5rem] px-5;

    &__title {
      @apply text-center text-xl font-semibold text-gray-800 dark:text-white sm:text-2xl;
    }

    &__icon {
      @apply mr-1;
    }

    &__description {
      @apply my-2 inline-block text-center text-sm font-semibold leading-5 text-gray-600 dark:text-slate-100 sm:text-base;
    }

    &__err-msg {
      @apply inline-block w-full text-center text-sm text-gray-400 dark:text-gray-200;
    }
  }
}
</style>
