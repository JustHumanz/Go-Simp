<template>
  <section class="ame-error">
    <img
      :src="`/assets/smolame/${img}.png`"
      alt=""
      class="ame-error-image"
    />
    <div class="ame-error-text">
      <h2 class="ame-error-text__title">
        <font-awesome-icon
          icon="circle-exclamation"
          class="fa-fw text-red-600 ml-2 ame-error-text__icon"
          v-if="type === `error`"
        />
        <font-awesome-icon
          icon="exclamation-triangle"
          class="fa-fw text-yellow-400 ml-2 ame-error-text__icon"
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
  @apply underline text-gray-700 hover:text-blue-400 dark:text-slate-200 dark:hover:text-slate-400 hover:decoration-2 underline-offset-1 decoration-dotted;
}
</style>

<style lang="scss" scoped>
.ame-error {
  @apply w-full h-screen flex flex-col md:flex-row items-center justify-center;

  &-image {
    @apply w-24 h-24 object-contain;
  }

  &-text {
    @apply max-w-[29.5rem] px-5;

    &__title {
      @apply text-center text-gray-800 dark:text-white font-semibold text-xl sm:text-2xl;
    }

    &__icon {
      @apply mr-1;
    }

    &__description {
      @apply inline-block text-center text-gray-600 dark:text-slate-100 font-semibold text-sm sm:text-base leading-5 my-2;
    }

    &__err-msg {
      @apply inline-block w-full text-center text-gray-400 dark:text-gray-200 text-sm;
    }
  }
}
</style>
