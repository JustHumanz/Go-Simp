<script setup>
import { RouterLink } from "vue-router"
</script>

<template>
  <div class="filter-nav">
    <div class="menu-filter">
      <a href="#" class="mobile-filter">Filters</a>
      <ul class="filters">
        <li class="filter">
          <a href="#" class="link-filter">Groups</a>
          <ul class="filter-submenu">
            <li class="pending" v-if="groups == null">
              <img src="/src/assets/amelia-watson-spin.gif" alt="" />
            </li>
            <li v-for="group in groups" :key="group._id">
              <router-link :to="`/vtuber/${group.ID}`">
                <img v-bind:src="group.GroupIcon" :alt="group.GroupName" />
                {{ group.GroupName }}
              </router-link>
            </li>
          </ul>
        </li>
        <li class="filter">
          <a href="#" class="link-filter">Regions</a>
          <ul class="filter-submenu">
            <li>Hi</li>
            <li>Hi</li>
          </ul>
        </li>
        <li class="filter">
          <a href="#" class="link-filter">More</a>
          <ul class="filter-submenu">
            <li>Hi</li>
          </ul>
        </li>
      </ul>
    </div>

    <div class="search">
      <input type="text" placeholder="Search..." />
    </div>
  </div>
  <section class="ame-page">
    <img src="/src/assets/amelia-watson-spin.gif" alt="" />
  </section>
  <!-- <section class="vtuber-page"></section> -->
</template>

<script>
import axios from "axios"
import Config from "../config.json"
export default {
  name: "Vtubers",
  data() {
    return {
      groups: null,
      group_id: -1,
      regions: [],
      vtubers: [],
    }
  },
  async mounted() {
    const data_groups = await axios
      .get(Config.REST_API + "/groups/")
      .then((response) => response.data)

    // add "all vtubers" in the first position of "groups"
    data_groups.unshift({
      ID: "",
      GroupName: "All Vtubers",
      GroupIcon: "https://i.imgur.com/XqQZQZL.png",
    })

    this.groups = data_groups
  },
  async created() {
    this.$watch(
      () => this.$route.params,
      async () => {
        await this.getRegions()
      },

      { immediate: true }
    )
  },
  beforeRouteEnter(to, from, next) {
    console.log(from)
    next()
  },
  methods: {
    async fetchVtubers() {
      if (this.group_id === this.$route.params.id || !this.$route.path.includes("vtuber")) {
        console.log("Clicked")
        return
      }

      this.vtubers = []

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

      await axios
        .get(Config.REST_API + "/members/", groupIdExist)
        .then((response) => {
          this.vtubers = response.data
        })
        .catch((error) => {
          console.log(error)
        })
      this.group_id = this.$route.params.id
      console.log(`Total data: ${this.vtubers.length}`)
    },

    async getRegions() {
      this.regions = []

      await this.fetchVtubers()

      await this.vtubers.forEach((vtuber) => {
        if (!this.regions.includes(vtuber.Region))
          this.regions.push(vtuber.Region)
      })

      console.log(this.regions)
    },
  },
}
</script>

<style lang="scss">
.filter-nav {
  @apply bg-blue-400 fixed top-16 w-screen py-3 px-3 md:px-[15%] lg:px-[17%] flex flex-wrap items-center sm:justify-between;

  .menu-filter {
    @apply flex-none;

    .mobile-filter {
      @apply sm:hidden;
    }

    &:focus-within .filters {
      @apply visible;
    }

    .filters {
      @apply flex absolute sm:relative flex-col sm:flex-row invisible sm:visible;

      .filter {
        @apply mx-1;
        .link-filter {
          @apply font-semibold text-white  px-2 py-1 rounded-md transition-shadow sm:hover:shadow-sm sm:hover:shadow-blue-600/75;
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

            li {
              a {
                @apply flex text-white w-full sm:pl-1 sm:pr-3 sm:py-1 sm:hover:bg-blue-500 font-semibold;

                img {
                  @apply w-6 h-auto drop-shadow-md mr-1;
                }
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
    @apply inline-block mx-1 ml-3 flex-auto sm:flex-none;

    input {
      @apply bg-blue-300 focus:bg-blue-200 py-1 px-2 rounded-md transition-shadow hover:shadow-sm hover:shadow-blue-600/75 focus:shadow-md focus:shadow-blue-600/75 w-full;
    }
  }
}

// make .ame-page full screen and placing content inside in to center
.ame-page {
  @apply flex flex-col items-center justify-center w-full h-screen;
}

.vtuber-page {
  @apply w-full md:w-[80%] lg:w-[75%] m-auto pt-[4.2rem];
}
</style>
