<template>
  <div class="filter-nav">
    <!--  -->
    <div class="menu-filter">
      <a href="#" class="mobile-filter"
        ><font-awesome-icon
          :icon="['fas', 'filter']"
          class="fa-fw fa-lg text-white"
          onclick="return false"
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
              <router-link :to="`/vtuber/${group.ID || ''}`" @click="resetLink">
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
                {{
                  (group.GroupName.charAt(0).toUpperCase() +
                  group.GroupName.slice(1)).replace("_", " ")
                }}
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
            <li
              class="filter-submenu-item"
              :class="{ active: !$route.query?.reg }"
              v-if="regions.length > 0"
            >
              <router-link
                :to="$route.href.replace(/\?reg=.{2}/, '')"
                @click="resetLink"
              >
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
              :class="{
                active:
                  $route.query?.reg &&
                  $route.query?.reg.toLowerCase() == region.code.toLowerCase(),
              }"
            >
              <router-link
                :to="
                  $route.href.replace(/\?reg=.{2}/, '') + `?reg=${region.code}`
                "
                @click="resetLink"
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

  <section
    class="vtuber-page"
    v-if="loaded && show_vtuber.length > 0 && !error_msg"
  >
    <div class="card-vtubers" v-for="vtuber in show_vtuber" :key="vtuber.ID">
      <div class="card-vtuber-img">
        <div class="tags">
          <!-- Get groupicon from group id -->
          <img
            class="group"
            v-if="vtuber.GroupID !== 10"
            :src="vtuber.Group.GroupIcon"
            :alt="vtuber.Group.GroupName"
          />
          <!-- Get flag icon from Region -->
          <img
            class="flag"
            v-bind:src="`/src/assets/flags/${vtuber.Regions.flagCode}.svg`"
            :alt="vtuber.Regions.name"
          />
          <a
            :href="vtuber.LiveURL"
            target="_blank"
            v-if="vtuber.IsBiliLive || vtuber.IsYtLive || vtuber.IsTwitchLive"
            class="live-indicator"
          >
            <!-- add icon youtube when live in youtube -->
            <font-awesome-icon
              v-if="vtuber.IsYtLive"
              :icon="['fab', 'youtube']"
              class="fa-fw"
            />

            <!-- add icon twitch when live in twitch -->
            <font-awesome-icon
              v-if="vtuber.IsTwitchLive"
              :icon="['fab', 'twitch']"
              class="fa-fw"
            />

            <!-- add icon bilibili when live in bilibili -->
            <font-awesome-icon
              v-if="vtuber.IsBiliLive"
              :icon="['fab', 'bilibili']"
              class="fa-fw"
            />
            LIVE
          </a>
          <div class="tooltip-info">
            {{
              `${vtuber.GroupID === 10 ? "Vtuber" : vtuber.Group.GroupName} ${
                vtuber.Regions.name
              }`
            }}
          </div>
        </div>
        <router-link :to="`/vtuber/members/${vtuber.ID}`">
          <div class="profile-pic" v-if="vtuber.Youtube">
            <img
              class="card-img-top"
              v-bind:src="vtuber.Youtube.Avatar.replace('s800', 's360')"
              onerror="this.src='/src/assets/smolame.jpg'"
              alt="Card image cap"
            />
          </div>
          <div
            class="profile-pic"
            v-else-if="vtuber.BiliBili && !vtuber.Youtube"
          >
            <img
              class="card-img-top"
              v-bind:src="`${vtuber.BiliBili.Avatar}@360w_360h_1c_1s.jpg`"
              referrerpolicy="no-referrer"
              onerror="this.src='/src/assets/smolame.jpg'"
              alt="Card image cap"
            />
          </div>
          <div v-else class="profile-pic">
            <img
              class="card-img-top"
              src="/src/assets/smolame.jpg"
              alt="Card image cap"
            />
          </div>
        </router-link>
        <div class="social-menu">
          <a
            v-if="vtuber.Youtube !== null"
            :href="`https://www.youtube.com/channel/${vtuber.Youtube.ID}`"
            target="_blank"
            class="social-icon youtube"
          >
            <font-awesome-icon :icon="['fab', 'youtube']" class="fa-fw" />
          </a>
          <a
            v-if="vtuber.BiliBili !== null"
            :href="`https://space.bilibili.com/${vtuber.BiliBili.ID}`"
            target="_blank"
            class="social-icon bilibili"
          >
            <font-awesome-icon :icon="['fab', 'bilibili']" class="fa-fw" />
          </a>
          <a
            v-if="vtuber.Twitch !== null"
            :href="`https://twitch.tv/${vtuber.Twitch.UserName}`"
            target="_blank"
            class="social-icon twitch"
          >
            <font-awesome-icon :icon="['fab', 'twitch']" class="fa-fw" />
          </a>
          <a
            v-if="vtuber.Twitter !== null"
            :href="`https://twitter.com/${vtuber.Twitter.UserName}`"
            target="_blank"
            class="social-icon twitter"
          >
            <font-awesome-icon :icon="['fab', 'twitter']" class="fa-fw" />
          </a>
        </div>
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
    v-if="
      !error_msg && (show_vtuber.length < getVtuberFilterData.length || !loaded)
    "
    :style="{
      height: loaded ? '' : '100vh',
    }"
  >
    <img
      :src="`/src/assets/loading/${Math.floor(Math.random() * 7)}.gif`"
      alt=""
    />
  </section>
  <!-- make section data vtuber not found -->
  <section class="error-page" v-if="loaded && show_vtuber.length < 1">
    <div class="error-image">
      <img src="/src/assets/smolame/bugs.png" alt="" />
    </div>
    <div class="error-text">
      <h2>
        <font-awesome-icon
          :icon="['fas', 'exclamation-triangle']"
          class="fa-fw"
        />
        <span>Hi, you cannot find your talent here</span>
      </h2>
      <p>
        The vtuber you are looking for is not found. Check group, region, or
        search. Or you can submit in
        <a
          href="https://github.com/JustHumanz/Go-Simp/issues/new?assignees=JustHumanz&labels=enhancement&template=add_vtuber.md&title=Add+%5BVtuber+Nickname%5D+from+%5BGroup%2FAgency%5D"
          target="_blank"
          rel="noopener noreferrer"
          >here</a
        >
      </p>
    </div>
  </section>
  <!-- make section when api is broken -->
  <section class="error-page" v-if="error_msg">
    <div class="error-image">
      <img src="/src/assets/smolame/lazer.png" alt="" />
    </div>
    <div class="error-text">
      <h2>
        <font-awesome-icon
          :icon="['fas', 'circle-exclamation']"
          class="fa-fw"
        />
        <span>Cannot call with server</span>
      </h2>
      <p>
        Check your internet connection, try restart your modem, or try again
        later.
      </p>
      <small>{{ error_msg }}</small>
    </div>
  </section>
</template>

<script>
import { RouterLink } from "vue-router"
import axios from "axios"
import Config from "../config.json"
import regionConfig from "../region.json"

import { library } from "@fortawesome/fontawesome-svg-core"
import {
  faYoutube,
  faBilibili,
  faTwitch,
  faTwitter,
} from "@fortawesome/free-brands-svg-icons"
import {
  faGlobeAmericas,
  faUser,
  faFilter,
  faMagnifyingGlass,
  faExclamationTriangle,
  faCircleExclamation,
} from "@fortawesome/free-solid-svg-icons"

library.add(
  faGlobeAmericas,
  faUser,
  faFilter,
  faMagnifyingGlass,
  faExclamationTriangle,
  faYoutube,
  faBilibili,
  faTwitch,
  faTwitter,
  faCircleExclamation
)

export default {
  name: "Vtubers",
  data() {
    return {
      loaded: false,
      error_msg: null,
      groups: [],
      groupID: -1,
      regions: [],
      vtubers: [],
      region_vtuber: [],
      show_vtuber: [],
      search: "",
      phName: "Search...",
      cancelGroups: null,
      cancelVtubers: null,
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
        this.search = ""

        if (
          this.groupID === this.$route.params?.id ||
          !this.$route.path.includes("vtuber") ||
          this.$route.path.includes("members") ||
          this.vtubers === []
        ) {
          console.log("Clicked")
          return
        }

        if (
          this.groups.length === 0 &&
          !this.loaded &&
          this.cancelGroups !== null
        )
          this.cancelGroups.cancel("Canceled")

        if (
          this.vtubers.length === 0 &&
          !this.loaded &&
          this.cancelVtubers !== null
        )
          this.cancelVtubers.cancel("Canceled")

        await this.getData()
      },

      { immediate: true }
    )
  },
  // beforeRouteEnter(to, from, next) {

  //   next()
  // },
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
      this.cancelVtubers = axios.CancelToken.source()

      const vtuber_data = await axios
        .get(Config.REST_API + "/members/", {
          cancelToken: this.cancelVtubers.token,
          ...groupIdExist,
        })
        .then((response) => {
          this.cancelVtubers = null
          return response.data
        })
        .catch((error) => {
          if (!axios.isCancel(error)) this.error_msg = error.message
        })

      if (this.error_msg) return false

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
    },

    async getGroupData() {
      if (this.groups.length > 0) return
      console.log("Fetching group data...")

      this.cancelGroups = axios.CancelToken.source()

      const data_groups = await axios
        .get(Config.REST_API + "/groups/", {
          cancelToken: this.cancelGroups.token,
        })
        .then((response) => response.data)
        .catch((error) => {
          if (!axios.isCancel(error)) this.error_msg = error.message
        })

      if (this.error_msg) return false

      // sort group data from GroupName
      data_groups.sort((a, b) => {
        if (a.GroupName < b.GroupName) return -1
        if (a.GroupName > b.GroupName) return 1
        return 0
      })

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
          Math.ceil(window.scrollY + window.innerHeight) >=
          document.body.offsetHeight
        let vtuber_show = this.show_vtuber.length

        if (bottomOfWindow && vtuber_show < this.getVtuberFilterData.length) {
          // count vtuber_show, then add 50 more
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
      this.loaded = false
      this.error_msg = null

      this.groupID = this.$route.params?.id

      // add title
      document.title = `Vtuber List`

      if ((await this.getGroupData()) === false) return
      if ((await this.fetchVtubers()) === false) return
      await this.getRegions()

      // add title "GroupName - Vtuber List" with first letter capital
      const NameGroup = this.groups.find(
        (group) => group.ID == this.groupID
      )?.GroupName
      document.title = `${
        NameGroup.charAt(0).toUpperCase() + NameGroup.slice(1)
      } - Vtuber List`

      this.loaded = true

      // limit vtubers from getVtuberData to 30
      this.phName =
        this.getVtuberFilterData[
          Math.floor(Math.random() * this.getVtuberFilterData.length)
        ]["EnName"]
      this.show_vtuber = this.getVtuberFilterData.slice(0, 30)
    },
    resetLink() {
      // get element .mobile-filter
      const mobileFilter = document.querySelector(".mobile-filter")
      mobileFilter.blur()

      // get element .filter-submenu-item a
      const filterSubmenuItem = document.querySelectorAll(
        ".filter-submenu-item a"
      )
      filterSubmenuItem.forEach((item) => {
        item.blur()
      })
    },
  },
}
</script>

<style lang="scss" scoped>
.filter-nav {
  @apply bg-blue-400 fixed top-16 py-3 px-4 w-screen flex flex-wrap items-center sm:justify-between md:justify-around justify-center z-10;

  .menu-filter {
    @apply flex-none;

    .mobile-filter {
      @apply sm:hidden;
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

        .filter-submenu {
          @apply hidden relative sm:absolute;

          .pending img {
            @apply w-16 h-16 object-cover sm:object-contain;
          }
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
                  @apply w-6 rounded-md object-contain drop-shadow-md mr-2 ml-1;
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
      }
    }
  }

  .search {
    @apply inline-block mx-1 ml-3 flex-auto sm:flex-none relative;

    .fa-magnifying-glass {
      @apply absolute mt-2 ml-2 text-blue-500;
    }

    input {
      @apply bg-blue-300 focus:bg-blue-200 disabled:bg-slate-500 py-1 px-2 rounded-lg transition-all hover:shadow-sm hover:shadow-blue-600/75 focus:shadow-md focus:shadow-blue-600/75 w-full text-gray-600 font-semibold placeholder:italic placeholder:text-blue-500 disabled:placeholder:text-blue-200 placeholder:font-normal pl-8 focus:outline-none;
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
  @apply w-full md:w-[80%] lg:w-[75%] p-4 m-auto pt-[4.2rem] mb-3 grid  gap-[1.5rem];
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));

  .card-vtubers {
    @apply bg-white shadow-md rounded-md overflow-hidden transition-transform;

    // scale up when hover using tailwind
    &:hover {
      @apply transform scale-105;
    }

    .card-vtuber-img {
      @apply relative;

      .tags {
        @apply absolute bg-gray-200/90 rounded-br-md h-6 flex items-center cursor-pointer;

        img {
          @apply object-contain w-5 ml-1 drop-shadow-md;

          &.flag {
            @apply mx-2 w-4 rounded-sm;
          }
        }
        .live-indicator {
          @apply bg-red-500/90 text-white px-2 py-1 rounded-br-md flex items-center text-xs font-semibold;

          svg {
            @apply mr-1;
          }
        }

        &:hover {
          .tooltip-info {
            @apply block;
          }
        }

        .tooltip-info {
          @apply absolute top-6 px-2 py-1 text-xs font-semibold text-gray-600 bg-gray-200/90 rounded-r-md z-[9] hidden whitespace-nowrap;
        }
      }

      .profile-pic {
        // size width and height square and fit to card
        @apply object-contain h-full aspect-square bg-smolame bg-cover;

        img {
          @apply object-contain h-full;
        }
      }
      .social-menu {
        @apply absolute bg-gray-200/90 bottom-0 right-0 p-1 rounded-tl-md;

        .social-icon {
          @apply p-2 text-gray-700;

          &.youtube:hover {
            // what color logo of youtube
            @apply text-youtube;
          }

          &.bilibili:hover {
            // what color logo of bilibili
            @apply text-bilibili;
          }

          &.twitch:hover {
            // what color logo of twitch
            @apply text-twitch;
          }

          &.twitter:hover {
            // what color logo of twitter
            @apply text-twitter;
          }
        }
      }
    }

    .card-vtuber-name {
      @apply p-3;

      h4 {
        @apply font-semibold text-base text-gray-800 leading-4;
      }

      small {
        @apply text-gray-600 leading-3;
      }
    }
  }
}

.error-page {
  @apply flex justify-center items-center w-full h-[95vh] flex-col sm:flex-row;

  .error-image {
    @apply w-24 h-24 object-contain;
  }

  // make .not-found-text fit with h2
  .error-text {
    @apply max-w-[29.5rem] px-5;

    h2 {
      @apply text-center text-gray-800 font-semibold text-xl sm:text-2xl;

      svg {
        @apply text-yellow-400 mr-1;
      }

      .fa-circle-exclamation {
        @apply text-red-500;
      }
    }

    p {
      @apply inline-block text-center text-gray-600 font-semibold text-sm sm:text-base leading-5 my-2;

      a {
        @apply text-gray-700 hover:text-blue-600 transition-all;
      }
    }
    small {
      @apply inline-block w-full text-center text-gray-400 text-sm;
    }
  }
}
</style>
