<script setup>
import ConfigView from "../components/DocsViews/Config.vue"
import TagingView from "../components/DocsViews/Taging.vue"
import UtilsView from "../components/DocsViews/Utils.vue"
</script>

<template>
  <div class="title" v-if="!error_msg">
    <span class="title__span">
      <font-awesome-icon
        :icon="['fas', 'circle-question']"
        class="title__svg"
      />
      Documentation
    </span>
  </div>
  <div class="container">
    <div class="tab">
      <ul>
        <li class="tab-list" :class="{ active: tab === 1 }">
          <a
            href="#"
            @click="tab = 1"
            class="tab-list__link"
            >Configuration</a
          >
        </li>
        <li class="tab-list" :class="{ active: tab === 2 }">
          <a
            href="#"
            @click="tab = 2"
            class="tab-list__link"
            >Roles and Taging</a
          >
        </li>
        <li class="tab-list" :class="{ active: tab === 3 }">
          <a
            href="#"
            @click="tab = 3"
            class="tab-list__link"
            >Utilites</a
          >
        </li>
        <li class="tab-list" :class="{ active: tab === 4 }">
          <a
            href="#"
            @click="tab = 4"
            class="tab-list__link"
            >More Command</a
          >
        </li>
      </ul>
    </div>
    <div class="message">
      There was a slight problem between the bot and the Discord API, but it
      should be resolved soon
    </div>
    <section v-if="tab === 1" class="content">
      <ConfigView />
    </section>
    <section v-if="tab === 2" class="content">
      <TagingView />
    </section>
    <section v-if="tab === 3" class="content">
      <UtilsView />
    </section>
  </div>
</template>

<script>
import { library } from "@fortawesome/fontawesome-svg-core"

// import {

// } from "@fortawesome/free-brands-svg-icons"
import { faCircleQuestion } from "@fortawesome/free-solid-svg-icons"

library.add(faCircleQuestion)

// change title
export default {
    data() {
        return {
            tab: 1,
        };
    },
    mounted() {
        document.title = "Documentation - Vtbot";
    },
    computed: {
        isActive() {
            return this.$route.path === "/docs";
        },
    },
    methods: {
        resetFocus() {
            this.$refs.search.focus();
        },
    },

}
</script>

<style lang="scss" scoped>
.title {
  @apply text-2xl font-semibold uppercase bg-blue-400 py-3 w-full flex flex-wrap;

  &__span {
    @apply w-[90%] md:w-[70%] lg:w-[65%] mx-auto text-white;
  }
  &__svg {
    @apply text-blue-200;
  }
}

.container {
  @apply mx-auto grid;
  grid-template-columns: 1fr 4fr;
  grid-template-rows: min-content 1fr;
  grid-template-areas:
    "tabs message"
    "tabs content";
}
.message {
  @apply bg-yellow-400 text-white px-4 py-2 mx-2 my-1 rounded-md font-semibold;
  grid-area: message;
}

.tab {
  @apply border-r-2 border-blue-400 flex-col py-2 hidden sm:inline-block;
  grid-area: tabs;

  &-list {
    @apply px-2 py-1 text-blue-400 hover:bg-blue-300;
    transition: all 0.2s ease-in-out;

    &.active {
      @apply bg-blue-300;
    }
  }

}

.content {
  @apply px-3;
  grid-area: content;
}
</style>
