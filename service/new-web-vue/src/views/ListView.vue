<script setup>
import NavbarList from "../components/MenuFilters/NavbarList.vue"
import AmeLoading from "../components/AmeComp/AmeLoading.vue"
import VtuberList from "../components/VtuberList.vue"
import AmeError from "../components/AmeComp/AmeError.vue"
</script>

<template>
  <NavbarList />
  <AmeLoading v-if="vtubersCount < 1 && !error_status" class="!h-screen" />
  <AmeError
    v-if="vtubersCount > 0 && filteredCount < 1"
    type="error"
    img="laptop"
    title="Your filter is incorrect"
    :description="`Check your filter. Or when your vtuber is not available, request ${link_request}`"
  />
  <AmeError
    v-if="vtubersCount > 0 && query && searchedCount < 1"
    type="warning"
    img="bugs"
    title="You find worng keyword"
    :description="`When your vtuber/member is not here, your can request ${link_request}, or try another keyword`"
  />
  <AmeError
    v-if="vtubersCount < 1 && error_status && error_status === 404"
    type="error"
    img="laptop"
    title="Your group is not available"
    :description="`Check another available group, or you can request a group ${link_request}`"
  />
  <AmeError
    v-if="vtubersCount < 1 && error_status && error_status !== 404"
    type="error"
    img="lazer"
    title="Something wrong when get request"
    :description="`Waiting for server response, restart your WiFi, or try again later`"
  />
  <VtuberList v-if="vtubersCount > 0" />
</template>

<script>
import { useMemberStore } from "@/stores/members.js"

export default {
  data() {
    return {
      group_id: null,
      // for menu handler
      activeListMenu: null,
      activeSubMenu: null,
      clickedSubMenu: false,
    }
  },
  async mounted() {
    await this.getVtuberData()

    // add abilty menu
    this.menuHandler()
  },
  computed: {
    link_request() {
      return `<a href="/new-vtuber" id="router-link" class="ame-error-text__link">here</a>`
    },
    filteredCount() {
      this.calculateFilters()
      return useMemberStore().members.filteredData.length
    },
    query() {
      return useMemberStore().members.query
    },
    searchedCount() {
      return useMemberStore().members.searchedData.length
    },
    vtubersCount() {
      return useMemberStore().members.data.length
    },
    error_status() {
      return useMemberStore().members.status
    },
  },
  methods: {
    async getVtuberData() {
      const router_before = window.history.state.back

      if (
        !router_before.includes("/vtuber") ||
        useMemberStore().members.group !== this.$route.params?.id
      )
        await useMemberStore().fetchMembers(this.$route.params?.id || null)
      useMemberStore().filterMembers()
      useMemberStore().sortingMembers()
    },

    async refreshNewData() {},

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
              this.activeSubMenu = null
              this.clickedSubMenu = false
              break
            default:
              this.activeListMenu = navbarFilter
              this.activeSubMenu = null
              this.clickedSubMenu = false
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
    async calculateFilters() {
      await new Promise((resolve) => setTimeout(resolve, 60))
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
