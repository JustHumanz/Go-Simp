<script setup>
import NavbarList from "../components/MenuFilters/NavbarList.vue"
import AmeLoading from "../components/AmeComp/AmeLoading.vue"
import VtuberList from "../components/VtuberList.vue"
import AmeError from "../components/AmeComp/AmeError.vue"
</script>

<template>
  <NavbarList
    :filters="filters"
    @search="getSearchData"
    :placeholder="phName"
    :disable_search="(null_data || !vtubers) && !search_query"
    :disabled="!vtubers"
  />
  <AmeLoading v-if="!vtubers && error_msg === ''" class="!h-screen" />
  <AmeError
    v-if="null_data && !search_query"
    type="error"
    img="laptop"
    title="Your filter is incorrect"
    :description="`Check your filter. Or when your vtuber is not available, request ${link_request}`"
  />
  <AmeError
    v-if="null_data && search_query"
    type="warning"
    img="bugs"
    title="You find worng keyword"
    :description="`When your vtuber/member is not here, your can request ${link_request}, or try another keyword`"
  />
  <AmeError
    v-if="!vtubers && error_msg === `Request failed with status code 404`"
    type="error"
    img="laptop"
    title="Your group is not available"
    :description="`Check another available group, or you can request a group ${link_request}`"
  />
  <AmeError
    v-if="!vtuber && error_msg !== '' && error_msg !== `Request failed with status code 404`"
    type="error"
    img="lazer"
    title="Something wrong when get request"
    :description="`Waiting for server response, restart your WiFi, or try again later`"
  />
  <VtuberList
    :vtubers="vtubers"
    :search_query="search_query"
    v-if="vtubers && vtubers.length > 0"
    @getPlaceholder="getPlaceholder"
    @null-data="nullData"
  />
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
      null_data: false,
      search_query: "",
      phName: "",
      error_msg: "",
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
        this.vtubers = null
        this.filters = null
        this.error_msg = ""
        this.getPlaceholder()
        console.log("Running...")

        await this.getVtuberData()
        if (this.error_msg) return
        this.getFilter()
      },
      { immediate: true }
    )
  },
  computed: {
    link_request() {
      return `<a href="https://github.com/JustHumanz/Go-Simp/issues/new?assignees=JustHumanz&labels=enhancement&template=add_vtuber.md&title=Add+%5BVtuber+Nickname%5D+from+%5BGroup%2FAgency%5D" target="_blank" class="ame-error-text__link">here</a>`
    },
  },
  methods: {
    async getVtuberData() {
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

      if (!vtuber_data) return

      vtuber_data.forEach((vtuber) => {
        regionConfig.forEach((region) => {
          if (region.code === vtuber.Region) {
            vtuber.Regions = region
          }
        })
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
    getSearchData(q) {
      this.search_query = q
    },

    getPlaceholder(name = "") {
      if (!this.vtubers || name == "") this.phName = "Search Vtubers..."
      else this.phName = name
    },

    nullData(bool) {
      this.null_data = bool
    },
  },
}
</script>
