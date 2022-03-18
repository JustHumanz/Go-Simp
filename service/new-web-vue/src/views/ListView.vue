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
    :error_status="error_status"
    :groups="groups"
  />
  <AmeLoading v-if="!vtubers && !error_status" class="!h-screen" />
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
    v-if="!vtubers && error_status && error_status === 404"
    type="error"
    img="laptop"
    title="Your group is not available"
    :description="`Check another available group, or you can request a group ${link_request}`"
  />
  <AmeError
    v-if="!vtubers && error_status && error_status !== 404"
    type="error"
    img="lazer"
    title="Something wrong when get request"
    :description="`Waiting for server response, restart your WiFi, or try again later`"
  />
  <VtuberList
    :vtubers="vtubers"
    :search_query="search_query"
    :groups="groups"
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
      groups: [],
      filters: null,
      group_id: null,
      error_status: null,
      null_data: false,
      search_query: "",
      phName: "",
    }
  },
  async created() {
    await this.getGroupData()

    this.$watch(
      () => this.$route.params,
      () => (this.group_id = this.$route.params?.id || null),
      { immediate: true }
    )

    this.$watch(
      () => this.group_id,
      async () => {
        if (!this.$route.path.includes("/vtubers")) return
        this.vtubers = null
        this.filters = null
        this.error_status = null
        this.getPlaceholder()
        console.log("Running...")

        await this.getVtuberData()
        if (this.error_status) return
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
          if (!axios.isCancel(error)) this.error_status = error.response.status
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

    async getGroupData() {
      if (this.groups.length > 0) return
      console.log("Fetching group data...")

      // this.cancelGroups = axios.CancelToken.source()

      const data_groups = await axios
        .get(Config.REST_API + "/v2/groups/", {
          // cancelToken: this.cancelGroups.token,
        })
        .then((response) => response.data)
        .catch((error) => {
          if (!axios.isCancel(error)) this.error_msg = error.message
        })

      if (this.error_msg) return false

      // sort group data from GroupName
      data_groups.sort((a, b) => {
        if (a.GroupName.toLowerCase() < b.GroupName.toLowerCase()) return -1
        if (a.GroupName.toLowerCase() > b.GroupName.toLowerCase()) return 1
        return 0
      })

      // add "all vtubers" in the first position of "groups"
      data_groups.unshift({
        GroupName: "All Vtubers",
        GroupIcon: "",
      })

      this.groups = data_groups
      console.log(`Total group: ${this.groups.length}`)
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
