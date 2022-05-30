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
    v-if="vtubers && vtubers.length > 0"
    @getPlaceholder="getPlaceholder"
    @null-data="nullData"
  />
</template>

<script>
import axios from "axios"
import Config from "../config.json"
import regionConfig from "../region.json"

import { useMemberStore } from "@/stores/members.js"

export default {
  data() {
    return {
      vtubers: null,
      filters: null,
      group_id: null,
      error_status: null,
      null_data: false,
      search_query: null,
      phName: "",
      // for menu handler
      activeListMenu: null,
      activeSubMenu: null,
      clickedSubMenu: false,
    }
  },
  async created() {
    console.log(window.location)

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
        await this.getFilter()
        this.calculateFilters()
      },
      { immediate: true }
    )

    // add abilty menu
    this.menuHandler()
  },
  computed: {
    link_request() {
      return `<a href="/new-vtuber" id="router-link" class="ame-error-text__link">here</a>`
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
    menuHandler() {
      // ON MOUSE DOWN
      window.addEventListener("mousedown", (e) => {
        const target = e.target.closest(".navbar-filter-item__link.sub-menu")

        if (!target) return
        if (!this.activeSubMenu || this.activeSubMenu !== target) {
          this.activeSubMenu = target
          this.clickedSubMenu = true
        }
      })

      // on mouse up
      window.addEventListener("mouseup", async (e) => {
        if (this.clickedSubMenu) {
          await new Promise((resolve) => setTimeout(resolve, 300))
          this.clickedSubMenu = false
        }
      })

      // on click
      document.addEventListener("click", (e) => {
        if (e.target.closest(".navbar-filter__link")) {
          const navbarFilter = e.target.closest(".navbar-filter__link")

          const liNavbarFilter = navbarFilter.parentElement

          if (liNavbarFilter.classList.contains("disabled")) {
            this.activeListMenu = null
            this.activeSubMenu = null
            this.clickedSubMenu = false
            navbarFilter.blur()
            return
          }

          switch (this.activeListMenu) {
            case navbarFilter:
              this.activeListMenu.blur()
              this.activeListMenu = null
              break
            case null:
              this.activeListMenu = navbarFilter
              break
            default:
              this.activeListMenu = navbarFilter
              break
          }
        } else if (e.target.closest(".navbar-filter-item__link")) {
          const navbarFilterItem = e.target.closest(".navbar-filter-item__link")

          if (!navbarFilterItem.classList.contains("sub-menu")) {
            this.activeListMenu = null
            navbarFilterItem.blur()
          } else {
            const filterSub = e.target.closest(
              ".navbar-filter-item__link.sub-menu"
            )

            if (this.activeSubMenu === filterSub && !this.clickedSubMenu) {
              console.log("closing submenu")
              this.activeSubMenu.blur()
              this.activeListMenu.focus()
              this.activeSubMenu = null
            }
          }
        } else if (e.target.closest(".navbar-submenu-item__link")) {
          const navbarSubItem = e.target.closest(".navbar-submenu-item__link")

          this.activeListMenu = null
          this.activeSubMenu = null
          navbarSubItem.blur()
        } else if (e.target.closest(".nav-search")) {
          if (this.activeListMenu !== null) {
            console.log("closing menu")
            this.activeListMenu = null
            this.activeSubMenu = null
          }

          const navbarSearchItem = e.target.closest(".nav-search")

          navbarSearchItem.children[1].focus()
        } else {
          if (this.activeListMenu === null) return
          console.log("closing menu")
          this.activeListMenu = null
          this.activeSubMenu = null
        }
      })

      // unfocus
      document.addEventListener("unfocus", (e) => {
        if (
          this.activeListMenu &&
          this.activeListMenu === document.activeElement
        )
          this.activeListMenu.blur()
        else if (
          !document.activeElement.classList.contains("nav-search__input")
        )
          document.activeElement.blur()

        this.activeListMenu = null
        this.activeSubMenu = null
      })
    },
    calculateFilters() {
      const subMenus = document.querySelectorAll(".navbar-submenu-items")

      // count children inside each submenu and substract 75
      subMenus.forEach((subMenu) => {
        const totalHeight = subMenu.children.length * 32
        // add --totalHeight in submenu
        subMenu.style.setProperty("--totalHeight", `${totalHeight}px`)
      })
    },
  },
}
</script>
