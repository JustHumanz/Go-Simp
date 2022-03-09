<script setup>
import AmeLoading from "../components/AmeComp/AmeLoading.vue"
import VtuberHeader from "../components/VtuberDetails/VtuberHeader.vue"
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
  <section class="vtuber-details" v-if="vtuber">
    <VtuberHeader :vtuber="vtuber" />
  </section>
</template>

<style lang="scss" scoped>
.title {
  @apply text-2xl font-semibold uppercase bg-blue-400 dark:bg-slate-500 py-2 w-full flex flex-wrap select-none fixed top-16 z-[9];

  &__span {
    @apply w-[95%] md:w-[75%] lg:w-[70%] mx-auto text-white;
  }
  &__svg {
    @apply text-blue-200 dark:text-gray-200;
  }
}

.back-button {
  @apply text-white hover:text-blue-300 dark:hover:text-gray-300 p-2;
}

.vtuber-details {
  @apply pb-4 w-full sm:w-[90%] md:w-[80%] lg:w-[75%] mx-auto mt-28;
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

    this.history_back = window.history.state.back || `/vtuber`

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
}
</script>
