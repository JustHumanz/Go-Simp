<script setup>
import AmeLoading from "../components/AmeComp/AmeLoading.vue"
import VtuberHeader from "../components/VtuberDetails/VtuberHeader.vue"
</script>

<template>
  <div class="title" v-if="!error_msg">
    <span class="title__span">
      <font-awesome-icon :icon="['fas', 'circle-info']" class="title__svg" />
      Info Vtuber
    </span>
  </div>
  <AmeLoading v-if="!vtuber && error_msg === ''" class="!h-screen" />
  <section class="vtuber-details">
    <VtuberHeader v-if="vtuber" :vtuber="vtuber" />
  </section>
</template>

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

.vtuber-details {
  @apply pb-4 w-full sm:w-[90%] md:w-[80%] lg:w-[75%] mx-auto;
}
</style>

<script>
import axios from "axios"
import Config from "../config.json"
import regionConfig from "../region.json"
import { CountTo } from "vue3-count-to"
import { library } from "@fortawesome/fontawesome-svg-core"

import { faCircleInfo } from "@fortawesome/free-solid-svg-icons"

library.add(faCircleInfo)

export default {
  data() {
    return {
      vtuber: null,
      fanart: null,
      error_msg: "",
    }
  },
  components: {
    CountTo,
  },
  async created() {
    document.title = "Vtuber Details - Vtbot"

    const member_data = await axios
      .get(Config.REST_API + "/v2/members/" + this.$route.params.id)
      .then((response) => response.data[0])
      .catch((error) => {
        this.error_msg = error.message
      })

    console.log(this.error_msg)

    if (this.error_msg) return

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
