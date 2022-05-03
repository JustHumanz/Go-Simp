<template>
  <header class="header">
    <img
      class="header-banner"
      draggable="false"
      :src="`${vtuber.Youtube.Banner.replace(
        's1200',
        ''
      )}s1707-fcrop64=1,00005a57ffffa5a8-k-c0xffffffff-no-nd-rj`"
      referrerpolicy="no-referrer"
      v-if="vtuber.Youtube"
      @error="onError"
      @contextmenu="disableContextMenu"
    />
    <img
      class="header-banner"
      draggable="false"
      :src="
        vtuber.BiliBili.Banner.replace('.jpg', '') + '@1707w_282h_1c_1s.jpg'
      "
      referrerpolicy="no-referrer"
      v-else-if="vtuber.BiliBili"
      @error="onError"
      @contextmenu="disableContextMenu"
    />
    <div class="header-banner" v-else @contextmenu="disableContextMenu" />
    <div class="header-info">
      <div
        class="header-profile-pic"
        :class="{ inactive: vtuber.Status == 'Inactive' }"
      >
        <img
          draggable="false"
          class="header-profile-pic__img"
          v-if="vtuber.Youtube"
          v-bind:src="vtuber.Youtube.Avatar.replace('s800', 's360')"
          referrerpolicy="no-referrer"
          onerror="this.src='/assets/smolame.jpg'"
          alt="Card image cap"
        />
        <img
          draggable="false"
          class="header-profile-pic__img"
          v-else-if="vtuber.BiliBili"
          v-bind:src="`${vtuber.BiliBili.Avatar}@360w_360h_1c_1s.jpg`"
          referrerpolicy="no-referrer"
          onerror="this.src='/assets/smolame.jpg'"
          alt="Card image cap"
        />
        <img
          draggable="false"
          class="header-profile-pic__img"
          v-else-if="vtuber.Twitch"
          v-bind:src="vtuber.Twitch.Avatar"
          referrerpolicy="no-referrer"
          onerror="this.src='/assets/smolame.jpg'"
          alt="Card image cap"
        />
        <img
          draggable="false"
          class="header-profile-pic__img"
          v-else
          src="/assets/smolame.jpg"
          alt="Card image cap"
        />
      </div>
      <div class="header-vtuber-name">
        <h4 class="header-vtuber-name__name">
          {{ vtuber.EnName }}{{ vtuber.JpName ? ` (${vtuber.JpName})` : "" }}
          <div class="nickname">
            {{ vtuber.NickName.toLowerCase() }}
            <span class="nickname-hover"
              >This nickname is used for initials when calling in the Discord
              command
              <router-link to="/docs/get-data-groups#get-vtuber-name"
                >Learn More</router-link
              ></span
            >
          </div>
        </h4>
        <div class="header-vtuber-name__group">
          <router-link
            :to="`/vtubers/${vtuber.Group.ID}`"
            class="header-vtuber-name__group-link"
          >
            <img
              :src="vtuber.Group.IconURL"
              :alt="vtuber.Group.GroupName"
              class="header-vtuber-name__group-icon"
              v-if="vtuber.Group.ID !== 10"
            />
            <img
              :src="`/assets/flags/${vtuber.Regions.flagCode}.svg`"
              :alt="vtuber.Group.GroupName"
              class="header-vtuber-name__group-icon"
              v-else
            />
            {{
              vtuber.Group.ID === 10
                ? "Vtuber"
                : vtuber.Group.GroupName.replace("_", " ")
            }}
            {{ vtuber.Regions.name }} </router-link
          ><router-link to="/vtubers?inac=true">
            {{ vtuber.Status === "Inactive" ? " (Inactive)" : "" }}
          </router-link>
          <div
            class="live-link"
            v-if="
              vtuber.IsLive.Youtube ||
              vtuber.IsLive.Twitch ||
              vtuber.IsLive.BiliBili
            "
          >
            <a
              :href="vtuber.IsLive.Youtube.URL"
              v-if="vtuber.IsLive.Youtube"
              target="_blank"
              rel="noopener noreferrer"
              class="live-link__link"
            >
              <font-awesome-icon
                :icon="['fab', 'youtube']"
                size="lg"
                class="live-link__link-icon"
              ></font-awesome-icon>
              <span class="live-link__link-text">LIVE</span>
            </a>
            <a
              :href="vtuber.IsLive.Twitch.URL"
              v-if="vtuber.IsLive.Twitch"
              target="_blank"
              rel="noopener noreferrer"
              class="live-link__link"
            >
              <font-awesome-icon
                :icon="['fab', 'twitch']"
                size="lg"
                class="live-link__link-icon"
              ></font-awesome-icon>
              <span class="live-link__link-text">LIVE</span>
            </a>
            <a
              :href="vtuber.IsLive.BiliBili.URL"
              v-if="vtuber.IsLive.BiliBili"
              target="_blank"
              rel="noopener noreferrer"
              class="live-link__link"
            >
              <font-awesome-icon
                :icon="['fab', 'bilibili']"
                size="lg"
                class="live-link__link-icon"
              ></font-awesome-icon>
              <span class="live-link__link-text">LIVE</span>
            </a>
          </div>
        </div>
      </div>
    </div>
  </header>
  <hr class="m-2 mt-3 mb-2" />
</template>

<script>
import { library } from "@fortawesome/fontawesome-svg-core"
import {
  faYoutube,
  faTwitch,
  faBilibili,
  faTwitter,
} from "@fortawesome/free-brands-svg-icons"

library.add(faYoutube, faTwitch, faBilibili, faTwitter)

export default {
  props: {
    vtuber: Object,
  },
  methods: {
    disableContextMenu(e) {
      e.preventDefault()
    },
    onError(e) {
      const parent = e.target.parentElement
      // delete img
      parent.removeChild(e.target)
      // create div with class .header-banner
      const banner = document.createElement("div")
      banner.classList.add("header-banner")
      // add scope id in banner attribute
      banner.setAttribute(this.$options.__scopeId, "")
      // append div to parent in first position
      parent.insertBefore(banner, parent.firstChild)
    },
  },
}
</script>

<style lang="scss" scoped>
.live-tag {
  @apply absolute bottom-0 inline-block w-full select-none bg-red-600 py-1 text-center text-xs font-semibold text-white;
}

.inactive {
  @apply bg-gray-600 bg-rip bg-contain bg-no-repeat;
  .header-profile-pic__img {
    @apply opacity-40 grayscale;
  }
}

.header {
  &-banner {
    @apply /*&h-[11.25rem]*/ /*bg-center bg-no-repeat*/ w-full bg-slate-200 bg-cover object-cover object-center dark:bg-slate-700;
    min-height: calc(16.1290322581vw - 1px);
    height: calc(16.1290322581vw - 1px);

    // @media (min-width: 640px) {
    //   height: calc(16.1290322581vw - 1px);
    // }

    @media (min-width: 768px) {
      min-height: calc((100vw - 240px) / 6.2 - 1px);
      height: calc((100vw - 240px) / 6.2 - 1px);
    }
  }
  &-info {
    @apply relative flex flex-col items-center sm:flex-row sm:pb-2;
  }
  &-profile-pic {
    @apply absolute -top-9 w-24 overflow-hidden rounded-md shadow-md sm:-top-14 sm:ml-10 sm:w-32;
  }
  &-vtuber-name {
    @apply mt-14 flex w-full flex-col items-center px-3 text-center sm:mt-0 sm:ml-[10.5rem] sm:items-stretch sm:text-left;

    &__name {
      @apply pt-3 text-xl font-semibold;
      .nickname {
        @apply relative block cursor-pointer text-xl font-thin text-gray-600 dark:text-slate-300 sm:inline-block;

        a {
          @apply text-white hover:text-gray-100 dark:hover:text-gray-300;
        }

        &-hover {
          @apply invisible absolute top-7 left-1/2 z-[5] inline-block w-44 -translate-x-1/2 -translate-y-5 rounded-md bg-sky-300 px-2 py-1 text-left text-sm font-normal opacity-0 delay-100 duration-300 ease-in-out dark:bg-slate-700;
          transition-property: opacity, transform;

          &::before {
            // add arrow up in center using tailwind
            @apply absolute -top-1 left-1/2 -translate-x-1/2 border-x-8 border-b-8 border-solid border-x-transparent border-b-sky-300 content-[''] dark:border-b-slate-700;
          }
        }

        &:hover {
          .nickname-hover {
            @apply visible translate-y-0 opacity-100;
          }
        }
      }
    }

    &__group {
      @apply mt-px inline-flex flex-col items-center text-sm font-light sm:flex-row;

      &-link {
        @apply text-gray-500 hover:text-gray-700 dark:text-gray-300 dark:hover:text-gray-400;
      }

      &-icon {
        @apply -mt-1 inline-block h-5 w-5 rounded-md object-contain;
      }
    }
  }
}

.live-link {
  @apply mt-2 flex w-fit space-x-1 text-xs sm:mt-0 sm:ml-2 sm:inline-block;

  &__link {
    @apply relative inline-flex -translate-y-px items-center space-x-1 rounded-full bg-red-600 py-1 px-2 font-semibold text-white shadow-sm shadow-red-700 transition duration-300 ease-in-out hover:translate-y-0 hover:text-white hover:shadow-none dark:shadow-red-300/25;
  }
}
</style>
