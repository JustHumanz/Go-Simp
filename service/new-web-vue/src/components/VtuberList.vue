<script setup>
import VtuberCard from "./VtuberCard.vue"
import AmeLoading from "./AmeComp/AmeLoading.vue"
</script>

<template>
  <a
    href="#"
    @click="scrollUp"
    onclick="return false"
    class="scroll-to-top"
    :style="{ transform: hide_scroll_up ? 'scale(0)' : 'scale(1)' }"
  >
    <font-awesome-icon icon="caret-up" class="fa-fw" />
  </a>

  <section v-if="$route.params?.id && group" class="group-detail">
    <img
      v-if="group.ID != 10"
      draggable="false"
      :src="group.GroupIcon"
      :alt="group.GroupName"
      class="group-detail__img"
    />
    <h4 class="group-detail__title">
      {{
        (
          group.GroupName.charAt(0).toUpperCase() + group.GroupName.slice(1)
        ).replace("_", " ")
      }}
      <span class="group-detail__command">
        {{ group.GroupName.toLowerCase() }}</span
      >
    </h4>
    <div class="group-detail__link" v-if="group.Youtube || group.ID === 10">
      <a
        :href="`https://www.youtube.com/channel/${youtube.YtChannel}`"
        target="_blank"
        rel="noopener noreferrer"
        v-for="youtube in group.Youtube"
        :key="youtube.Region"
        class="group-detail__link-item"
      >
        <font-awesome-icon :icon="['fab', 'youtube']" class="fa-fw" />
        <span>{{ regions.find((r) => r.code === youtube.Region).name }}</span>
      </a>
      <span v-if="group.ID === 10" class="text-gray-600 dark:text-gray-200">
        This is a District Group,
        <router-link
          to="/docs/quick-start#about-independent-group"
          class="text-gray-700 hover:text-gray-800 dark:text-gray-50 dark:hover:text-gray-300"
        >
          Read more</router-link
        >
      </span>
    </div>
  </section>
  <section class="vtuber-list" v-if="!nullData">
    <VtuberCard
      v-for="vtuber in limitedVtubers"
      :key="vtuber.ID"
      :vtuber="vtuber"
    />
  </section>

  <AmeLoading v-if="limitedVtubers.length < vtubers.length" class="my-4" />
</template>

<script>
import { library } from "@fortawesome/fontawesome-svg-core"
import { faCaretUp } from "@fortawesome/free-solid-svg-icons"
import regionConfig from "../regions.json"

library.add(faCaretUp)

import { useGroupStore } from "@/stores/groups"
import { useMemberStore } from "@/stores/members.js"

export default {
  emits: ["getPlaceholder", "null-data"],
  data() {
    return {
      limitedVtubers: [],
      nullData: false,
      hide_scroll_up: true,
      regions: null,
    }
  },
  created() {
    this.regions = regionConfig
    const store = useMemberStore()

    this.$watch(
      () => this.vtubers || store.members.query,
      () => {
        this.limitedVtubers = this.vtubers.slice(0, 30)
      },
      { immediate: true }
    )

    this.ScrollFuncions()
  },
  computed: {
    group() {
      return (
        useGroupStore().groups.data.find(
          (group) => group.ID == this.$route.params?.id
        ) || null
      )
    },
    vtubers() {
      const store = useMemberStore()

      return store.members.searchedData.length > 0
        ? store.members.searchedData
        : store.members.filteredData
    },
    test() {
      return this.vtubers.length
    },
  },
  methods: {
    ScrollFuncions() {
      window.onscroll = () => {
        let bottomOfWindow =
          Math.ceil(window.scrollY + window.innerHeight) >=
          document.body.offsetHeight - 145

        let vtubers_count = this.limitedVtubers.length

        if (bottomOfWindow && vtubers_count < this.vtubers.length) {
          // count vtubers_count, then add 25 more
          console.log("Extend more data...")
          this.limitedVtubers = this.vtubers.slice(0, vtubers_count + 25)
        }

        let topOfWindow = window.scrollY <= 0

        if (topOfWindow) {
          this.hide_scroll_up = true
        } else {
          this.hide_scroll_up = false
        }
      }
    },
    scrollUp() {
      window.scrollTo({
        top: 0,
        behavior: "smooth",
      })
    },
  },
}
</script>

<style lang="scss" scoped>
.vtuber-list {
  @apply mx-auto grid w-[95%] gap-[1.5rem] pb-4 sm:w-[85%] md:w-[80%] lg:w-[75%];
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
}

.group-detail {
  @apply mx-auto mb-3 grid w-[95%] gap-1 rounded-md bg-slate-200 p-2 shadow-md dark:bg-slate-500 sm:w-[85%] md:w-[80%] lg:w-[75%];
  grid-template-columns: max-content auto;
  grid-template-areas: "image title" "image link";

  &__img {
    @apply m-1 h-14 w-14 self-center justify-self-start object-cover drop-shadow-md;
    grid-area: image;
  }

  &__title {
    @apply ml-1 self-center justify-self-start text-center text-2xl font-bold text-gray-700 dark:text-white;
    grid-area: title;

    // when &__title is end item
    &:last-child {
      grid-area: 1 / 2 / 3 / 3;
    }
  }

  &__command {
    @apply text-xl font-thin text-slate-800 dark:text-slate-200;
  }

  &__link {
    @apply ml-1;
    grid-area: link;

    &-item {
      @apply m-[2px] inline-flex items-center space-x-2 rounded-full bg-youtube px-3 py-1 text-white hover:brightness-90;
    }
  }
}

.scroll-to-top {
  @apply fixed bottom-0 right-0 z-[2] m-4 rounded-full bg-sky-400 px-3 pb-2 pt-3
  text-2xl text-white shadow-sm shadow-sky-700/50 transition-transform hover:shadow-md dark:bg-slate-800 dark:shadow-slate-100/50 sm:mr-[2.5%] md:mr-[7.5%] lg:mr-[10.5%];
}
</style>
