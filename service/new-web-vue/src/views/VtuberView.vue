<template>
  <div class="filter-nav">
    <div class="menu-filter">
      <a href="#" class="mobile-filter"
        ><font-awesome-icon
          :icon="['fas', 'filter']"
          class="fa-fw fa-lg text-white"
      /></a>
      <ul class="filters">
        <li class="filter">
          <a href="#" class="link-filter" onclick="return false">Groups</a>
          <ul class="filter-submenu">
            <li class="pending" v-if="groups == null">
              <img
                :src="`/src/assets/loading/${Math.floor(
                  Math.random() * 7
                )}.gif`"
                alt=""
              />
            </li>
            <li
              v-for="group in groups"
              :key="group.ID"
              class="filter-submenu-item"
              :class="{ active: groupID == group.ID }"
            >
              <router-link :to="`/vtuber/${group.ID || ''}`">
                <img
                  v-if="group.GroupIcon"
                  :src="group.GroupIcon"
                  :alt="group.GroupName"
                />
                <font-awesome-icon
                  v-else
                  :icon="['fas', 'user']"
                  class="fa-fw"
                />
                {{ group.GroupName }}
              </router-link>
            </li>
          </ul>
        </li>
        <li class="filter">
          <a href="#" class="link-filter" onclick="return false">Regions</a>
          <ul class="filter-submenu">
            <li class="pending" v-if="regions.length < 1">
              <img
                :src="`/src/assets/loading/${Math.floor(
                  Math.random() * 7
                )}.gif`"
                alt=""
              />
            </li>
            <li class="filter-submenu-item" :class="{active: !$route.query?.reg}" v-if="regions.length > 0">
              <router-link :to="$route.href.replace(/\?reg=.{2}/, '')">
                <font-awesome-icon
                  :icon="['fas', 'earth-americas']"
                  class="fa-fw"
                />
                All Regions
              </router-link>
            </li>
            <li
              v-for="region in regions"
              :key="region._id"
              class="filter-submenu-item"
              :class="{active: $route.query?.reg && $route.query?.reg.toLowerCase()  == region.code.toLowerCase()}"
            >
              <router-link
                :to="
                  $route.href.replace(/\?reg=.{2}/, '') + `?reg=${region.code}`
                "
              >
                <img
                  v-if="region.flagCode"
                  :src="`/src/assets/flags/${region.flagCode}.svg`"
                  alt=""
                />
                <font-awesome-icon
                  v-else
                  :icon="['fas', 'earth-americas']"
                  class="fa-fw"
                />
                {{ region.name }}
              </router-link>
            </li>
          </ul>
        </li>
      </ul>
    </div>

    <div class="search">
      <font-awesome-icon
        :icon="['fas', 'magnifying-glass']"
        class="fa-fw fa-md text-white"
      />
      <input
        type="text"
        v-model="search"
        id="vtuber-search"
        :placeholder="phName"
        :disabled="vtubers.length < 1"
      />
    </div>
  </div>
  <section class="vtuber-page" :class="{hide: !loaded}">
    <div class="card-vtubers" v-for="vtuber in show_vtuber" :key="vtuber.ID">
      <div class="card-vtuber-img">
        <div class="tags">
          <!-- Get groupicon from group id -->
          <img
            v-if="vtuber.GroupID !== 10"
            :src="vtuber.Group.GroupIcon"
            :alt="vtuber.Group.GroupName"
          />
          <!-- Get flag icon from Region -->
          <img
            v-bind:src="`/src/assets/flags/${vtuber.Regions.flagCode}.svg`"
            :alt="vtuber.Regions.name"
          />
          <a
            :href="vtuber.LiveURL"
            target="_blank"
            :class="{
              hide:
                !vtuber.IsBiliLive && !vtuber.IsYtLive && !vtuber.IsTwitchLive,
            }"
            class="live-indicator"
          >
            LIVE
          </a>
          <div class="tooltip-info">
            {{ `${vtuber.Group.GroupName} ${vtuber.Region}` }}
          </div>
        </div>
        <router-link :to="`/vtuber/members/${vtuber.ID}`">
          <div class="profile-pic" v-if="vtuber.Youtube !== null">
            <img
              class="card-img-top"
              v-bind:src="vtuber.Youtube.Avatar"
              onerror="this.src='/src/assets/smolame.png'"
              alt="Card image cap"
            />
          </div>
          <div class="profile-pic" v-else-if="vtuber.BiliBili !== null">
            <img
              class="card-img-top"
              v-bind:src="vtuber.BiliBili.Avatar"
              onerror="this.src='/src/assets/smolame.png'"
              alt="Card image cap"
            />
          </div>
          <div v-else class="profile-pic">
            <img
              class="card-img-top"
              src="/src/assets/smolame.png"
              alt="Card image cap"
            />
          </div>
        </router-link>
      </div>
      <router-link :to="`/vtuber/members/${vtuber.ID}`">
        <div class="card-vtuber-name">
          <h4>{{ vtuber.EnName }}</h4>
          <small>{{ vtuber.NickName }}</small>
        </div>
      </router-link>
    </div>
  </section>
  <section
    class="ame-page"
    :class="{
      hide: show_vtuber >= getVtuberFilterData && loaded,
      hiscreen: !loaded,
    }"
  >
    <img
      :src="`/src/assets/loading/${Math.floor(Math.random() * 7)}.gif`"
      alt=""
    />
  </section>
</template>

<script>
import { RouterLink } from "vue-router"
import axios from "axios"
import Config from "../config.json"
import regionConfig from "../region.json"

import { library } from "@fortawesome/fontawesome-svg-core"
// import { fa } from '@fortawesome/free-brands-svg-icons'
import {
  faGlobeAmericas,
  faUser,
  faFilter,
  faMagnifyingGlass,
} from "@fortawesome/free-solid-svg-icons"

library.add(faGlobeAmericas, faUser, faFilter, faMagnifyingGlass)

export default {
  name: "Vtubers",
  data() {
    return {
      loaded: false,
      groups: null,
      groupID: -1,
      regions: [],
      vtubers: [],
      region_vtuber: [],
      show_vtuber: [],
      search: "",
      phName: "Search...",
    }
  },
  async mounted() {
    this.Searching()
    this.changingRegion()
    this.ExtendVtuberData()
  },
  async created() {
    this.$watch(
      () => this.$route.params,
      async () => {
        if (
          this.groupID === this.$route.params?.id ||
          !this.$route.path.includes("vtuber") ||
          this.vtubers === []
        ) {
          console.log("Clicked")
          return
        }

        await this.getData()
      },

      { immediate: true }
    )
  },
  beforeRouteEnter(to, from, next) {
    console.log(from)
    next()
  },
  computed: {
    getVtuberFilterData() {
      // get reg from ?reg=
      const reg = this.$route.query.reg
      const vtuber_data = reg
        ? this.vtubers.filter(
            (vtuber) => vtuber.Region.toLowerCase() === reg.toLowerCase()
          )
        : this.vtubers

      console.log("Filtering")

      return vtuber_data.filter((post) => {
        let EnName = post.EnName.toLowerCase().includes(
          this.search.toLowerCase()
        )
        let JpName
        if (post.JpName != null) {
          JpName = post.JpName.toLowerCase().includes(this.search.toLowerCase())
        }
        return EnName || JpName
      })
    },
  },
  methods: {
    async fetchVtubers() {
      this.vtubers = []
      this.regions = []

      const groupIdExist = this.$route.params.id
        ? {
            crossDomain: true,
            params: {
              groupid: this.$route.params.id,
              live: "true",
            },
          }
        : {}

      // List vtuber
      console.log("Fetching data...")

      this.locate = this.$route.href
      this.loaded = false

      const vtuber_data = await axios
        .get(Config.REST_API + "/members/", groupIdExist)
        .then((response) => response.data)

      console.log("Add more stuff...")

      vtuber_data.forEach((vtuber) => {
        this.groups.forEach((group) => {
          if (group.ID === vtuber.GroupID) vtuber.Group = group
        })

        regionConfig.forEach((region) => {
          if (region.code === vtuber.Region) {
            vtuber.Regions = region
          }
        })
      })

      // sort vtuber_data from nameen
      vtuber_data.sort((a, b) => {
        if (a.EnName < b.EnName) return -1
        if (a.EnName > b.EnName) return 1
        return 0
      })

      // sort vtuber_data when IsYtLive true or IsBiliLive true or IsTwitchLive true first
      vtuber_data.sort((a, b) => {
        if (a.IsYtLive && !b.IsYtLive) return -1
        if (!a.IsYtLive && b.IsYtLive) return 1
        if (a.IsBiliLive && !b.IsBiliLive) return -1
        if (!a.IsBiliLive && b.IsBiliLive) return 1
        if (a.IsTwitchLive && !b.IsTwitchLive) return -1
        if (!a.IsTwitchLive && b.IsTwitchLive) return 1
        return 0
      })

      this.vtubers = vtuber_data

      console.log(`Total vtuber: ${this.vtubers.length}`)
      this.loaded = true
    },

    async getGroupData() {
      if (this.groups !== null) return
      console.log("Fetching group data...")

      const data_groups = await axios
        .get(Config.REST_API + "/groups/")
        .then((response) => response.data)

      // add "all vtubers" in the first position of "groups"
      data_groups.unshift({
        GroupName: "All Vtubers",
        GroupIcon: "",
      })

      // change GroupIcon from ID 10 to ""
      data_groups.forEach((group) => {
        if (group.ID === 10) group.GroupIcon = ""
      })

      this.groups = data_groups
      console.log(`Total group: ${this.groups.length}`)
    },

    async getRegions() {
      console.log("Get regions...")

      const region_data = []

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
      })

      console.log(`Total region: ${region_data.length}`)

      // sort region_data from name
      region_data.sort((a, b) => {
        if (a.name < b.name) return -1
        if (a.name > b.name) return 1
        return 0
      })

      this.regions = region_data

      console.log(this.regions)
    },
    ExtendVtuberData() {
      window.onscroll = () => {
        let bottomOfWindow =
          document.documentElement.scrollTop + window.innerHeight ===
          document.documentElement.offsetHeight

        if (bottomOfWindow) {
          // count vtuber_show, then add 50 more
          let vtuber_show = this.show_vtuber.length
          this.show_vtuber = this.getVtuberFilterData.slice(0, vtuber_show + 30)
        }
      }
    },
    Searching() {
      // Check this.search changing
      console.log("Searching...")
      this.$watch(
        () => this.search,
        () => {
          this.show_vtuber = this.getVtuberFilterData.slice(0, 30)
        },
        { immediate: true }
      )
    },
    changingRegion() {
      this.$watch(
        () => this.$route.query.reg,
        () => {
          if (this.getVtuberFilterData.length > 0)
            this.phName =
              this.getVtuberFilterData[
                Math.floor(Math.random() * this.getVtuberFilterData.length)
              ]["EnName"]
          this.show_vtuber = this.getVtuberFilterData.slice(0, 30)
        },
        { immediate: true }
      )
    },
    async getData() {
      this.groupID = this.$route.params?.id || null

      await this.getGroupData()
      await this.fetchVtubers()
      await this.getRegions()

      // limit vtubers from getVtuberData to 30
      this.phName =
        this.getVtuberFilterData[
          Math.floor(Math.random() * this.getVtuberFilterData.length)
        ]["EnName"]
      this.show_vtuber = this.getVtuberFilterData.slice(0, 30)
    },
  },
}
</script>

<style lang="scss">
.filter-nav {
  @apply bg-blue-400 fixed top-16 w-screen py-3 px-3 md:px-[15%] lg:px-[17%] flex flex-wrap items-center sm:justify-between z-10;

  .menu-filter {
    @apply flex-none;

    .mobile-filter {
      @apply sm:hidden;

      svg > path {
        // @apply fill-white;
      }
    }

    &:focus-within .filters {
      @apply visible;
    }

    .filters {
      @apply flex absolute sm:relative flex-col sm:flex-row invisible sm:visible bg-blue-400 sm:bg-transparent left-0 sm:left-auto mt-4 sm:m-0 w-screen sm:w-auto max-h-[86.5vh] sm:max-h-[none] overflow-y-auto overflow-x-hidden sm:overflow-visible;

      .filter {
        @apply sm:mx-1;
        .link-filter {
          @apply font-semibold text-white px-3 py-1 sm:px-2 md:rounded-md transition-shadow
          w-full h-full inline-block hover:bg-blue-500/50 sm:hover:bg-transparent sm:hover:shadow-sm sm:hover:shadow-blue-600/75;
        }

        &:focus-within {
          .link-filter {
            @apply sm:shadow-md sm:shadow-blue-600;
          }

          .filter-submenu {
            @apply block bg-blue-400 sm:shadow-center sm:shadow-blue-600/75 sm:mt-1 sm:rounded-md transition-all sm:overflow-y-auto sm:overflow-x-hidden sm:max-h-60;
            @media (min-width: 640px) {
              scrollbar-width: none; /* Firefox */
              -ms-overflow-style: none; /* IE 10+ */
              &::-webkit-scrollbar {
                /* Chromium and Safari */
                display: none;
              }
            }

            // add event for li exept .pending

            .filter-submenu-item {
              a {
                @apply flex text-white w-full px-7 sm:pl-1 sm:pr-3 py-1 hover:bg-blue-500/50 font-semibold sm:w-44;

                img {
                  @apply w-6 object-contain drop-shadow-md mr-2 ml-1;
                }
                svg {
                  @apply py-1 px-[2px] mr-2 ml-1;
                }
              }
              &.active {
                @apply bg-blue-500;
              }
            }
          }
        }

        .filter-submenu {
          @apply hidden relative sm:absolute;
        }
      }
    }
  }

  .search {
    @apply inline-block mx-1 ml-3 flex-auto sm:flex-none relative;

    .fa-magnifying-glass {
      @apply absolute mt-2 ml-2 text-blue-500;
    }

    input {
      @apply bg-blue-300 focus:bg-blue-200 py-1 px-2 rounded-lg transition-shadow hover:shadow-sm hover:shadow-blue-600/75 focus:shadow-md focus:shadow-blue-600/75 w-full text-gray-600 font-semibold placeholder:italic placeholder:text-blue-500 placeholder:font-normal pl-8;
    }
  }
}

// make .ame-page full screen and placing content inside in to center
.ame-page {
  @apply flex flex-col items-center justify-center w-full h-40;

  img {
    @apply h-32 object-contain;
  }
}

.vtuber-page {
  @apply w-full md:w-[80%] lg:w-[75%] p-4 m-auto pt-[4.2rem] mb-3 grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-[1.5rem];

  .card-vtubers {
    @apply bg-white shadow-md rounded-md overflow-hidden transition-transform;

    // scale up when hover using tailwind
    &:hover {
      @apply transform scale-105;
    }

    .card-vtuber-img {
      @apply relative;

      .tags {
        @apply absolute bg-gray-200 rounded-br-md h-8 flex items-center cursor-pointer;

        img {
          @apply object-contain w-8 mx-1 drop-shadow-md;
        }
        .live-indicator {
          @apply bg-red-500 text-white px-2 py-1 rounded-br-md;
        }

        &:hover {
          .tooltip-info {
            @apply block;
          }
        }

        .tooltip-info {
          @apply absolute top-8 px-2 py-1 text-xs font-semibold text-gray-600 bg-gray-200 rounded-r-md z-[9] hidden whitespace-nowrap;
        }
      }

      .profile-pic {
        // size width and height square and fit to card
        @apply object-contain h-full aspect-square bg-smolame bg-cover;
        img {
          @apply object-contain h-full;
        }
      }
    }

    .card-vtuber-name {
      @apply px-4 py-2 mt-2;

      h4 {
        @apply font-semibold text-lg text-gray-800 leading-4;
      }

      small {
        @apply text-gray-600 leading-3;
      }
    }
  }
}

.hide {
  @apply hidden;
}

.hiscreen {
  @apply h-screen;
}
</style>
