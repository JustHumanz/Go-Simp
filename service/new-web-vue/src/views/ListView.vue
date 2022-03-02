<script setup>
import NavbarList from "../components/MenuFilters/NavbarList.vue"
import AmeLoading from "../components/AmeComp/AmeLoading.vue"
import VtuberList from "../components/VtuberList.vue"
</script>

<template>
  <NavbarList :filters="filters" />
  <AmeLoading v-if="!vtubers" class="!h-screen" />
  <VtuberList :vtubers="vtubers" v-if="vtubers" />
</template>

<script>
import axios from "axios"
import Config from "../config.json"
import regionConfig from "../region.json"

export default {
  data() {
    return {
      vtubers: null,
      filters: null,
      group_id: null,
    }
  },
  async created() {
    this.$watch(
      () => this.$route.params,
      () => (this.group_id = this.$route.params?.id || null),
      { immediate: true }
    )

    this.$watch(
      () => this.group_id,
      async () => {
        console.log("Running...")
        this.filters = null

        await this.getVtuberData()
        this.getFilter()
      },
      { immediate: true }
    )
  },
  computed: {
    FilteredVtuber() {
      
    },
  },
  methods: {
    async getVtuberData() {
      this.vtubers = null

      const groupId = this.$route.params.id
        ? { groupid: this.$route.params.id }
        : {}

      const groupIdExist = {
        params: {
          ...groupId,
          live: "true",
        },
      }

      const vtuber_data = await axios
        .get(Config.REST_API + "/v2/members/", {
          // cancelToken: this.cancelVtubers.token,
          ...groupIdExist,
        })
        .then((response) => {
          // this.cancelVtubers = null
          return response.data
        })
        .catch((error) => {
          if (!axios.isCancel(error)) this.error_msg = error.message
        })

      this.vtubers = vtuber_data
      console.log("Total vtuber members: " + this.vtubers.length)
    },
    async getFilter() {
      // Get region
      const region_data = []

      let twitch = false
      let youtube = false
      let twitter = false
      let bilibili = false
      let inactive = false

      await this.vtubers.forEach((vtuber) => {
        // loop regionConfig
        regionConfig.forEach((region) => {
          if (region.code === vtuber.Region) {
            // check if region already exist
            if (region_data.find((region) => region.code === vtuber.Region)) {
              return
            }

            region_data.push({
              code: vtuber.Region,
              name: region.name,
              flagCode: region.flagCode,
            })
          }
        })

        // Sort region_data from A to Z with toLowerCase
        region_data.sort((a, b) => {
          if (a.name.toLowerCase() < b.name.toLowerCase()) return -1
          if (a.name.toLowerCase() > b.name.toLowerCase()) return 1
          return 0
        })

        if (vtuber.BiliBili) bilibili = true
        if (vtuber.Twitter) twitter = true
        if (vtuber.Youtube) youtube = true
        if (vtuber.Twitch) twitch = true

        if (vtuber.Status == "Inactive") inactive = true
      })

      console.log("Total region: " + region_data.length)

      this.filters = {
        region: region_data,
        bilibili,
        twitter,
        youtube,
        twitch,
        inactive,
      }
    },
  },
}
</script>
