<template>
  <div class="title">
    <span class="title__span">
      <font-awesome-icon icon="circle-plus" class="title__svg fa-fw" />
      Add new vtuber
    </span>
  </div>
  <div class="form">
    <transition
      enter-active-class="slide-active"
      enter-from-class="slide-right"
      enter-to-class="slide-center"
      leave-active-class="slide-active"
      leave-from-class="slide-center"
      leave-to-class="slide-left"
    >
      <div v-if="step === 1">
        Lorem ipsum dolor sit amet consectetur adipisicing elit. Quae inventore
        architecto explicabo in quibusdam possimus totam ut temporibus.
        Architecto beatae facere odit aspernatur quos eveniet! Consequuntur
        doloribus qui et eveniet at eligendi maxime, nihil id distinctio nulla
        neque laudantium ipsum voluptate ad aliquid. Reprehenderit consequuntur
        ea velit ex quam autem, deserunt, nulla asperiores, hic quae eaque.
        Eaque illum inventore reprehenderit recusandae fugit culpa ut qui
        maiores dolor, sapiente, excepturi ea, maxime expedita. Voluptatem
        blanditiis qui obcaecati dolorum quisquam, ratione optio nesciunt
        dolorem, consequatur voluptatibus hic? Sint, non. Autem id voluptatem
        asperiores quos dignissimos odio sequi quasi aperiam sed inventore.
      </div>
    </transition>
    <transition></transition>
  </div>
  <button @click="showHide">hide/show</button>
</template>

<script>
import { library } from "@fortawesome/fontawesome-svg-core"
import { faCirclePlus } from "@fortawesome/free-solid-svg-icons"

library.add(faCirclePlus)

export default {
  data() {
    return {
      step: 1,
    }
  },
  methods: {
    showHide(e) {
      this.step = this.step === 3 ? 1 : this.step + 1
      console.log(this.show)
    },
  },
}
</script>

<style lang="scss" scoped>
.title {
  @apply fixed top-16 z-[9] flex w-full select-none flex-wrap bg-blue-400 py-2 text-2xl font-semibold uppercase dark:bg-slate-500;

  &__span {
    @apply mx-auto w-[95%] text-white md:w-[75%] lg:w-[70%];
  }
  &__svg {
    @apply text-blue-200 dark:text-gray-200;
  }
}

.form {
  @apply mx-auto mt-28 w-full overflow-x-hidden sm:w-[90%] md:w-[80%] lg:w-[75%];
}

.slide {
  &-left {
    @apply -translate-x-[110%];
  }

  &-center {
    @apply translate-x-0;
  }

  &-right {
    @apply translate-x-[110%];
  }

  &-active {
    @apply transition-transform duration-300 ease-in-out;
  }
}
</style>
