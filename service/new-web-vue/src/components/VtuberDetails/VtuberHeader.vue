<template>
  <header class="header">
    <img
      class="header-banner bg-cyan-300"
      draggable="false"
      :src="`${vtuber.Youtube.Banner.replace(
        's1200',
        ''
      )}s1707-fcrop64=1,00005a57ffffa5a8-k-c0xffffffff-no-nd-rj`"
      referrerpolicy="no-referrer"
      v-if="vtuber.Youtube"
      @contextmenu="disableContextMenu"
    />
    <img
      class="header-banner bg-cyan-300"
      draggable="false"
      :src="
        vtuber.BiliBili.Banner.replace('.jpg', '') + '@1707w_282h_1c_1s.jpg'
      "
      referrerpolicy="no-referrer"
      v-else-if="vtuber.BiliBili"
      @contextmenu="disableContextMenu"
    />
    <div
      class="header-banner bg-cyan-300"
      v-else
      @contextmenu="disableContextMenu"
    />
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
          onerror="this.src='/src/assets/smolame.jpg'"
          alt="Card image cap"
        />
        <img
          draggable="false"
          class="header-profile-pic__img"
          v-else-if="vtuber.BiliBili"
          v-bind:src="`${vtuber.BiliBili.Avatar}@360w_360h_1c_1s.jpg`"
          referrerpolicy="no-referrer"
          onerror="this.src='/src/assets/smolame.jpg'"
          alt="Card image cap"
        />
        <img
          draggable="false"
          class="header-profile-pic__img"
          v-else-if="vtuber.Twitch"
          v-bind:src="vtuber.Twitch.Avatar"
          referrerpolicy="no-referrer"
          onerror="this.src='/src/assets/smolame.jpg'"
          alt="Card image cap"
        />
        <img
          draggable="false"
          class="header-profile-pic__img"
          v-else
          src="/src/assets/smolame.jpg"
          alt="Card image cap"
        />
        <span
          class="live-tag"
          v-if="
            vtuber.IsLive.Youtube ||
            vtuber.IsLive.Twitch ||
            vtuber.IsLive.BiliBili
          "
        >
          <font-awesome-icon
            :icon="['fab', 'youtube']"
            class="fa-fw"
            v-if="vtuber.IsLive.Youtube"
          />
          <font-awesome-icon
            :icon="['fab', 'twitch']"
            class="fa-fw"
            v-if="vtuber.IsLive.Twitch"
          />
          <font-awesome-icon
            :icon="['fab', 'bilibili']"
            class="fa-fw"
            v-if="vtuber.IsLive.BiliBili"
          />
          LIVE
        </span>
      </div>
      <div class="header-vtuber-name">
        <h4 class="header-vtuber-name__name">
          {{ vtuber.EnName }}{{ vtuber.JpName ? ` (${vtuber.JpName})` : "" }}
          <div class="nickname">
            {{ vtuber.NickName }}
            <span class="nickname-hover"
              >This nickname is used for initials when calling in the Discord
              command
              <router-link to="/docs#members">Learn More</router-link></span
            >
          </div>
        </h4>
        <div class="header-vtuber-name__group">
          <router-link :to="`/vtubers/${vtuber.Group.ID}`">
            <img
              :src="vtuber.Group.IconURL"
              :alt="vtuber.Group.GroupName"
              class="header-vtuber-name__group-icon"
              v-if="vtuber.Group.ID !== 10"
            />
            <img
              :src="`/src/assets/flags/${vtuber.Regions.flagCode}.svg`"
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
        </div>
      </div>
    </div>
  </header>
  <hr class="m-2 mt-3 mb-2" />
  <div class="link-header">
    <a
      :href="`https://youtube.com/channel/${vtuber.Youtube.YoutubeID}`"
      v-if="vtuber.Youtube"
      target="_blank"
      rel="noopener noreferrer"
      class="link-header__link bg-youtube"
    >
      <font-awesome-icon
        :icon="['fab', 'youtube']"
        size="lg"
        class="link-header__link-icon"
      ></font-awesome-icon>
      <span class="link-header__link-text">YouTube</span>
    </a>
    <a
      :href="`https://twitch.tv/${vtuber.Twitch.Username}`"
      v-if="vtuber.Twitch"
      target="_blank"
      rel="noopener noreferrer"
      class="link-header__link bg-twitch"
    >
      <font-awesome-icon
        :icon="['fab', 'twitch']"
        size="lg"
        class="link-header__link-icon"
      ></font-awesome-icon>
      <span class="link-header__link-text">Twitch</span>
    </a>
    <a
      :href="`https://live.bilibili.com/${vtuber.BiliBili.LiveID}`"
      v-if="vtuber.BiliBili"
      target="_blank"
      rel="noopener noreferrer"
      class="link-header__link bg-bilibili"
    >
      <font-awesome-icon
        :icon="['fab', 'bilibili']"
        size="lg"
        class="link-header__link-icon"
      ></font-awesome-icon>
      <span class="link-header__link-text">Live</span>
    </a>
    <a
      :href="`https://space.bilibili.com/${vtuber.BiliBili.SpaceID}`"
      v-if="vtuber.BiliBili"
      target="_blank"
      rel="noopener noreferrer"
      class="link-header__link bg-bilibili"
    >
      <font-awesome-icon
        :icon="['fab', 'bilibili']"
        size="lg"
        class="link-header__link-icon"
      ></font-awesome-icon>
      <span class="link-header__link-text">Space</span>
    </a>
    <a
      class="link-header__link bg-twitter"
      :href="`https://twitter.com/${vtuber.Twitter.Username}`"
      v-if="vtuber.Twitter"
      target="_blank"
      rel="noopener noreferrer"
    >
      <font-awesome-icon
        :icon="['fab', 'twitter']"
        size="lg"
        class="link-header__link-icon"
      ></font-awesome-icon>
      <span class="link-header__link-text">@{{ vtuber.Twitter.Username }}</span>
    </a>
    <a
      :href="`https://twitter.com/hashtag/${vtuber.Twitter.Fanart.replace(
        '#',
        ''
      )}`"
      v-if="vtuber.Twitter && vtuber.Twitter.Fanart"
      target="_blank"
      rel="noopener noreferrer"
      class="link-header__link bg-twitter"
    >
      <font-awesome-icon
        icon="paint-brush"
        size="lg"
        class="link-header__link-icon"
      ></font-awesome-icon>
      <span class="link-header__link-text"
        >Twitter (#{{ vtuber.Twitter.Fanart.replace("#", "") }})</span
      >
    </a>
    <a
      :href="`https://www.pixiv.net/en/tags/${fanart_pixiv}`"
      target="_blank"
      rel="noopener noreferrer"
      class="link-header__link bg-pixiv"
    >
      <font-awesome-icon
        icon="paint-brush"
        size="lg"
        class="link-header__link-icon"
      />
      <span class="link-header__link-text"> Pixiv (#{{ fanart_pixiv }}) </span>
    </a>
  </div>
  <hr class="m-2" />
  <div
    class="link-header"
    v-if="
      vtuber.IsLive.Youtube || vtuber.IsLive.Twitch || vtuber.IsLive.BiliBili
    "
  >
    <a
      :href="vtuber.IsLive.Youtube.URL"
      v-if="vtuber.IsLive.Youtube"
      target="_blank"
      rel="noopener noreferrer"
      class="link-header__link bg-red-600 hover:bg-red-700"
    >
      <font-awesome-icon
        :icon="['fab', 'youtube']"
        size="lg"
        class="link-header__link-icon"
      ></font-awesome-icon>
      <span class="link-header__link-text">LIVE on YouTube</span>
    </a>
    <a
      :href="vtuber.IsLive.Twitch.URL"
      v-if="vtuber.IsLive.Twitch"
      target="_blank"
      rel="noopener noreferrer"
      class="link-header__link bg-red-600 hover:bg-red-700"
    >
      <font-awesome-icon
        :icon="['fab', 'twitch']"
        size="lg"
        class="link-header__link-icon"
      ></font-awesome-icon>
      <span class="link-header__link-text">LIVE on Twitch</span>
    </a>
    <a
      :href="vtuber.IsLive.BiliBili.URL"
      v-if="vtuber.IsLive.BiliBili"
      target="_blank"
      rel="noopener noreferrer"
      class="link-header__link bg-red-600 hover:bg-red-700"
    >
      <font-awesome-icon
        :icon="['fab', 'bilibili']"
        size="lg"
        class="link-header__link-icon"
      ></font-awesome-icon>
      <span class="link-header__link-text">LIVE on BiliBili</span>
    </a>
  </div>
    <hr
    class="m-2 mt-3 mb-2"
    v-if="
      vtuber.IsLive.Youtube || vtuber.IsLive.Twitch || vtuber.IsLive.BiliBili
    "
  />
  
</template>

<script>
import { library } from "@fortawesome/fontawesome-svg-core"
import {
  faYoutube,
  faTwitch,
  faBilibili,
  faTwitter,
} from "@fortawesome/free-brands-svg-icons"
import { faPaintBrush } from "@fortawesome/free-solid-svg-icons"

library.add(faYoutube, faTwitch, faBilibili, faTwitter, faPaintBrush)

export default {
  props: {
    vtuber: Object,
  },

  computed: {
    fanart_pixiv() {
      return this.vtuber.JpName
        ? this.vtuber.JpName.split("/")[0]
            .split(" ")
            .reduce((acc, cur) => (acc + acc !== "" ? "・" : "" + cur), "")
        : this.vtuber.EnName.split(" ").reduce((acc, cur) => acc + cur, "")
    },
  },
  methods: {
    disableContextMenu(e) {
      e.preventDefault()
    },
  },
}
</script>

<style lang="scss" scoped>
.live-tag {
  @apply w-full inline-block bg-red-600 text-white font-semibold text-center text-xs py-1 absolute bottom-0 select-none;
}

.inactive {
  @apply bg-rip bg-contain bg-no-repeat bg-gray-600;
  .header-profile-pic__img {
    @apply grayscale opacity-40;
  }
}

.header {
  &-banner {
    @apply w-full /*&h-[11.25rem]*/ /*bg-center bg-cover bg-no-repeat*/ object-cover object-center;
    height: calc(16.1290322581vw - 1px);

    // @media (min-width: 640px) {
    //   height: calc(16.1290322581vw - 1px);
    // }

    @media (min-width: 768px) {
      height: calc((100vw - 240px) / 6.2 - 1px);
    }
  }
  &-info {
    @apply flex items-center flex-col sm:flex-row sm:pb-2 relative;
  }
  &-profile-pic {
    @apply absolute -top-9 sm:-top-14 sm:ml-10 w-24 sm:w-32 rounded-md shadow-md overflow-hidden;
  }
  &-vtuber-name {
    @apply w-full flex items-center sm:items-stretch mt-14 px-3 sm:mt-0 sm:ml-[10.5rem] flex-col text-center sm:text-left;

    &__name {
      @apply text-xl font-semibold pt-3;
      .nickname {
        @apply text-gray-600 dark:text-slate-300 text-xl font-thin block sm:inline-block relative cursor-pointer;

        a {
          @apply text-white hover:text-gray-100 dark:hover:text-gray-300;
        }

        &-hover {
          @apply bg-sky-300 dark:bg-slate-700 text-sm font-normal px-2 py-1 rounded-md w-44 text-left inline-block absolute top-7 left-1/2 -translate-x-1/2 -translate-y-5 invisible opacity-0 delay-100 duration-300 ease-in-out;
          transition-property: opacity, transform;

          &::before {
            // add arrow up in center using tailwind
            @apply content-[''] border-x-8 border-b-8 border-solid border-x-transparent border-b-sky-300 dark:border-b-slate-700 absolute -top-1 left-1/2 -translate-x-1/2;
          }
        }

        &:hover {
          .nickname-hover {
            @apply opacity-100 translate-y-0 visible;
          }
        }
      }
    }

    &__group {
      @apply text-sm font-light text-gray-500 dark:text-gray-300 inline-block mt-px;

      a {
        @apply hover:text-gray-700 dark:hover:text-gray-400;
      }

      &-icon {
        @apply w-5 h-5 object-contain rounded-md inline-block -mt-1;
      }
    }
  }
}

.link-header {
  @apply flex justify-center items-center px-2 flex-wrap;

  &__link {
    @apply inline-flex items-center text-white px-3 m-px py-2 rounded-full;

    &-icon {
      @apply w-5 h-5 object-contain rounded-md mr-2;
    }
    &-text {
      @apply text-sm font-semibold;
    }

    &:hover {
      @apply brightness-90;
    }
  }

  & > span {
    @apply hidden;
  }
}
</style>
