<script setup>
import AmeLoading from "../components/AmeComp/AmeLoading.vue"
import AmeError from "../components/AmeComp/AmeError.vue"
import VtuberHeader from "../components/VtuberDetails/VtuberHeader.vue"
import Detail from "../components/VtuberDetails/Detail.vue"
import YoutubeCount from "../components/VtuberDetails/YoutubeCount.vue"
import TwitchCount from "../components/VtuberDetails/TwitchCount.vue"
import BiliBiliCount from "../components/VtuberDetails/BiliBiliCount.vue"
import TwitterCount from "../components/VtuberDetails/TwitterCount.vue"
</script>

<template>
  <div class="title" v-if="!error_status">
    <span class="title__span">
      <router-link class="back-button" :to="history_back">
        <font-awesome-icon icon="caret-left" class="fa-fw" />
      </router-link>
      <font-awesome-icon icon="circle-info" class="title__svg fa-fw" />
      Info Vtuber
    </span>
  </div>
  <AmeLoading v-if="!vtuber && !error_status" class="!h-screen" />
  <AmeError
    v-if="!vtuber && error_status && error_status === 404"
    type="error"
    img="laptop"
    title="Your vtuber is not available"
    :description="`Check another available vtuber, or you can request a member vtuber ${link_request}`"
  />
  <section class="vtuber-details" v-if="vtuber">
    <VtuberHeader :vtuber="vtuber" />
    <Detail :vtuber="vtuber" />
    <YoutubeCount :youtube="vtuber.Youtube" v-if="vtuber.Youtube" />
    <TwitchCount :twitch="vtuber.Twitch" v-if="vtuber.Twitch" />
    <BiliBiliCount :bilibili="vtuber.BiliBili" v-if="vtuber.BiliBili" />
    <TwitterCount :twitter="vtuber.Twitter" v-if="vtuber.Twitter" />
  </section>
</template>

<style lang="scss" scoped>
.title {
  @apply sticky top-16 z-[9] flex w-full select-none flex-wrap bg-blue-400 py-2 text-2xl font-semibold uppercase dark:bg-slate-500;

  &__span {
    @apply mx-auto w-[95%] text-white md:w-[75%] lg:w-[70%];
  }
  &__svg {
    @apply text-blue-200 dark:text-gray-200;
  }
}

.back-button {
  @apply p-2 text-white hover:text-blue-300 dark:hover:text-gray-300;
}

.vtuber-details {
  @apply mx-auto w-full pb-4 sm:w-[90%] md:w-[80%] lg:w-[75%];
}
</style>

<script>
import axios from "axios"
import Config from "../config.json"
import regionConfig from "../region.json"
import { CountTo } from "vue3-count-to"
import { library } from "@fortawesome/fontawesome-svg-core"

import { faCircleInfo, faCaretLeft } from "@fortawesome/free-solid-svg-icons"

library.add(faCircleInfo, faCaretLeft)

export default {
  data() {
    return {
      vtuber: null,
      fanart: null,
      error_status: "",
      history_back: null,
    }
  },
  components: {
    CountTo,
  },
  async created() {
    document.title = "Vtuber Details - Vtbot"

    this.history_back = window.history.state.back || `/vtubers`

    const member_data = await axios
      .get(Config.REST_API + "/v2/members/" + this.$route.params.id)
      .then((response) => response.data[0])
      .catch((error) => {
        this.error_status = error.response.status
      })

    if (this.error_status) return

    regionConfig.forEach((region) => {
      if (region.code === member_data.Region) {
        member_data.Regions = region
      }
    })

    this.vtuber = member_data
    console.log(this.vtuber)

    document.title = this.vtuber.EnName + " - Vtuber Details"
  },
  computed: {
    link_request() {
      return `<a href="/new-vtuber" id="router-link" class="ame-error-text__link">here</a>`
    },
  },
}
</script>
