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

  <div class="pt-24 xs:pt-14" />

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

  <AmeLoading
    v-if="limitedVtubers.length < searchVtubers.length"
    class="my-4"
  />
</template>

<script>
import { library } from "@fortawesome/fontawesome-svg-core"
import { faCaretUp } from "@fortawesome/free-solid-svg-icons"
import regionConfig from "../region.json"

library.add(faCaretUp)

import { useGroupStore } from "@/stores/groups"

export default {
  props: {
    vtubers: {
      type: Array,
    },
    search_query: {
      type: String,
      default: null,
    },
  },
  emits: ["getPlaceholder", "null-data"],
  data() {
    return {
      limitedVtubers: [],
      nullData: false,
      hide_scroll_up: true,
      regions: null,
    }
  },
  async created() {
    this.regions = regionConfig

    this.$watch(
      () => this.$route.query,
      async () => {
        this.limitedVtubers = await this.searchVtubers.slice(0, 30)
      },
      { immediate: true }
    )

    this.$watch(
      () => this.search_query,
      async () => {
        this.limitedVtubers = await this.searchVtubers.slice(0, 30)
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
    filteredVtubers() {
      let vtuber_data = this.vtubers

      if (!this.$route.path.includes("/vtubers")) return []

      // Filter vtuber.Region from this.$route.query.reg
      vtuber_data = vtuber_data.filter((vtuber) => {
        return this.$route.query.reg
          ? vtuber.Region == this.$route.query.reg
          : vtuber
      })

      // Filter platform inside vtuber

      switch (this.$route.query.plat) {
        case "yt":
          // Filter when vtuber_data.Youtube is not null
          vtuber_data = vtuber_data.filter((vtuber) => {
            return vtuber.Youtube != null
          })
          break
        case "tw":
          // Filter when vtuber_data.Twitch is not null
          vtuber_data = vtuber_data.filter((vtuber) => {
            return vtuber.Twitch != null
          })
          break
        case "bl":
          // Filter when vtuber_data.BiliBili is not null
          vtuber_data = vtuber_data.filter((vtuber) => {
            return vtuber.BiliBili != null
          })
          break
      }

      // Filter when vtuber_data.Status is "Inactive"

      vtuber_data = vtuber_data.filter((vtuber) => {
        return this.$route.query.inac == "true"
          ? vtuber.Status === "Inactive"
          : vtuber
      })

      console.log("Total after filtered: " + vtuber_data.length)

      this.$emit(
        "getPlaceholder",
        vtuber_data.length > 0
          ? vtuber_data[Math.floor(Math.random() * vtuber_data.length)][
              "EnName"
            ]
          : ""
      )

      return vtuber_data
    },

    sortingVtubers() {
      // get query.sort when exist
      let vtuber_data = this.filteredVtubers

      if (!this.$route.path.includes("/vtubers")) return []

      // Sorting vtuber by EnName DESC and lowercase first
      vtuber_data = vtuber_data.sort((a, b) => {
        if (a.EnName.toLowerCase() > b.EnName.toLowerCase()) return -1
        if (a.EnName.toLowerCase() < b.EnName.toLowerCase()) return 1
        return 0
      })

      switch (this.$route.query.sort) {
        case undefined:
          console.log("Sort by Alphabet")
          // Sorting vtuber by EnName ASC and lowercase
          vtuber_data = vtuber_data.sort((a, b) => {
            if (a.EnName.toLowerCase() < b.EnName.toLowerCase()) return -1
            if (a.EnName.toLowerCase() > b.EnName.toLowerCase()) return 1
            return 0
          })
          break
        case "-name":
          console.log("Sort by Reverse Alphabet")
          break
        case "yt":
          console.log("Sort by Most Youtube Subscriber")
          // Sorting vtuber by Youtube.Subscriber ASC when exist
          vtuber_data = vtuber_data.sort((a, b) => {
            if (a.Youtube) {
              if (a.Youtube?.Subscriber < b.Youtube?.Subscriber) return 1
              if (a.Youtube?.Subscriber > b.Youtube?.Subscriber) return -1
            } else return 1
            return 0
          })
          break
        case "-yt":
          console.log("Sort by Least Youtube Subscriber")
          // Sorting vtuber by Youtube.Subscriber DESC when exist
          vtuber_data = vtuber_data.sort((a, b) => {
            if (b.Youtube) {
              if (a.Youtube?.Subscriber > b.Youtube?.Subscriber) return 1
              if (a.Youtube?.Subscriber < b.Youtube?.Subscriber) return -1
            } else return 1
            return 0
          })
          break
        case "tw":
          console.log("Sort by Most Twitch Followers")
          //Sorting vtuber by Twitch.Followers ASC when exist
          vtuber_data = vtuber_data.sort((a, b) => {
            if (a.Twitch) {
              if (a.Twitch?.Followers < b.Twitch?.Followers) return 1
              if (a.Twitch?.Followers > b.Twitch?.Followers) return -1
            } else return 1
            return 0
          })
          break
        case "-tw":
          console.log("Sort by Least Twitch Followers")
          //Sorting vtuber by Twitch.Followers DESC when exist
          vtuber_data = vtuber_data.sort((a, b) => {
            if (b.Twitch) {
              if (a.Twitch?.Followers > b.Twitch?.Followers) return 1
              if (a.Twitch?.Followers < b.Twitch?.Followers) return -1
            } else return 1
            return 0
          })
          break
        case "bl":
          console.log("Sort by Most Bilibili Followers")
          //Sorting vtuber by BiliBili.Followers ASC when exist
          vtuber_data = vtuber_data.sort((a, b) => {
            if (a.BiliBili) {
              if (a.BiliBili?.Followers < b.BiliBili?.Followers) return 1
              if (a.BiliBili?.Followers > b.BiliBili?.Followers) return -1
            } else return 1
            return 0
          })
          break
        case "-bl":
          console.log("Sort by Least Bilibili Followers")
          //Sorting vtuber by BiliBili.Followers DESC when exist
          vtuber_data = vtuber_data.sort((a, b) => {
            if (b.BiliBili) {
              if (a.BiliBili?.Followers > b.BiliBili?.Followers) return 1
              if (a.BiliBili?.Followers < b.BiliBili?.Followers) return -1
            } else return 1
            return 0
          })
          break
        case "twr":
          console.log("Sort by Most Twitter Followers")
          //Sorting vtuber by Twitter.Followers ASC when exist
          vtuber_data = vtuber_data.sort((a, b) => {
            if (a.Twitter) {
              if (a.Twitter?.Followers < b.Twitter?.Followers) return 1
              if (a.Twitter?.Followers > b.Twitter?.Followers) return -1
            } else return 1
            return 0
          })
          break
        case "-twr":
          console.log("Sort by Least Twitter Followers")
          //Sorting vtuber by Twitter.Followers DESC when exist
          vtuber_data = vtuber_data.sort((a, b) => {
            if (b.Twitter) {
              if (a.Twitter?.Followers > b.Twitter?.Followers) return 1
              if (a.Twitter?.Followers < b.Twitter?.Followers) return -1
            } else return 1
            return 0
          })
          break
      }

      // Sorting vtuber when IsLive.BiliBili is not null (Object), IsLive.Twitch is not null (Object), and IsLive.Youtube is not null (Object)
      return vtuber_data.sort((a, b) => {
        if (!a.IsLive.Youtube && b.IsLive.Youtube) return 1
        if (a.IsLive.Youtube && !b.IsLive.Youtube) return -1
        if (!a.IsLive.Twitch && b.IsLive.Twitch) return 1
        if (a.IsLive.Twitch && !b.IsLive.Twitch) return -1
        if (!a.IsLive.BiliBili && b.IsLive.BiliBili) return 1
        if (a.IsLive.BiliBili && !b.IsLive.BiliBili) return -1
        return 0
      })
    },

    searchVtubers() {
      // Filter vtuber.EnName or vtuber.JpName from this.search_query
      const vtuber_data = this.sortingVtubers.filter((post) => {
        let EnName = post.EnName.toLowerCase().includes(
          this.search_query ? this.search_query.toLowerCase() : ""
        )
        let JpName
        if (post.JpName != null) {
          JpName = post.JpName.toLowerCase().includes(
            this.search_query ? this.search_query.toLowerCase() : ""
          )
        }
        return EnName || JpName
      })

      this.nullData = vtuber_data.length === 0
      this.$emit("null-data", this.nullData)
      return vtuber_data
    },
  },
  methods: {
    ScrollFuncions() {
      window.onscroll = () => {
        let bottomOfWindow =
          Math.ceil(window.scrollY + window.innerHeight) >=
          document.body.offsetHeight - 145

        let vtubers_count = this.limitedVtubers.length

        if (bottomOfWindow && vtubers_count < this.searchVtubers.length) {
          // count vtubers_count, then add 25 more
          console.log("Extend more data...")
          this.limitedVtubers = this.searchVtubers.slice(0, vtubers_count + 25)
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
  @apply mx-auto mb-3 grid w-[95%] gap-1 rounded-md bg-blue-300 p-2 shadow-md dark:bg-slate-500 sm:w-[85%] md:w-[80%] lg:w-[75%];
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
