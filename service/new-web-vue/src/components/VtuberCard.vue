<template>
  <div class="card-vtuber" :class="{ inactive: vtuber.Status == 'Inactive' }">
    <div class="card-vtuber-image">
      <div class="tag-vtuber">
        <div class="tag-vtuber-agency">
          <img
            draggable="false"
            class="tag-vtuber-agency__icon"
            :src="vtuber.Group.IconURL"
            :alt="vtuber.Group.GroupName"
            v-if="vtuber.Group.ID !== 10"
          />
          <img
            draggable="false"
            class="tag-vtuber-agency__flag"
            :src="`/src/assets/flags/${vtuber.Regions.flagCode}.svg`"
            :alt="vtuber.Regions.name"
          />
          <font-awesome-icon
            class="fa-fw"
            icon="skull"
            v-if="vtuber.Status == 'Inactive'"
          />
          <span class="tag-vtuber-agency__hover">{{
            `${vtuber.Group.ID === 10 ? `Vtuber` : GroupName} ${
              vtuber.Regions.name
            } ${vtuber.Status === "Inactive" ? ` (Inactive)` : ""}`
          }}</span>
        </div>
        <a
          class="tag-vtuber-live"
          v-if="
            vtuber.IsLive.Youtube ||
            vtuber.IsLive.Twitch ||
            vtuber.IsLive.BiliBili
          "
          :href="liveLink"
          target="_blank"
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
          <span>LIVE</span>
        </a>
      </div>
      <router-link :to="`/vtuber/${vtuber.ID}`" class="card-vtuber-image__link">
        <img
          draggable="false"
          class="card-vtuber-image__img"
          v-if="vtuber.Youtube"
          v-bind:src="vtuber.Youtube.Avatar.replace('s800', 's360')"
          referrerpolicy="no-referrer"
          onerror="this.src='/src/assets/smolame.jpg'"
          alt="Card image cap"
        />
        <img
          draggable="false"
          class="card-vtuber-image__img"
          v-else-if="vtuber.BiliBili"
          v-bind:src="`${vtuber.BiliBili.Avatar}@360w_360h_1c_1s.jpg`"
          referrerpolicy="no-referrer"
          onerror="this.src='/src/assets/smolame.jpg'"
          alt="Card image cap"
        />
        <img
          draggable="false"
          class="card-vtuber-image__img"
          v-else-if="vtuber.Twitch"
          v-bind:src="vtuber.Twitch.Avatar"
          referrerpolicy="no-referrer"
          onerror="this.src='/src/assets/smolame.jpg'"
          alt="Card image cap"
        />
        <img
          draggable="false"
          class="card-vtuber-image__img"
          v-else
          src="/src/assets/smolame.jpg"
          alt="Card image cap"
        />
      </router-link>
      <div class="vtuber-link">
        <a
          :href="`https://youtube.com/channel/${vtuber.Youtube.YoutubeID}`"
          target="_blank"
          class="vtuber-link__link hover:!text-youtube"
          rel="noopener noreferrer"
          v-if="vtuber.Youtube"
        >
          <font-awesome-icon
            :icon="['fab', 'youtube']"
            class="fa-fw"
            rel="noopener noreferrer"
            v-if="vtuber.Youtube"
          />
        </a>
        <a
          :href="`https://twitch.tv/${vtuber.Twitch.Username}`"
          target="_blank"
          class="vtuber-link__link hover:!text-twitch"
          rel="noopener noreferrer"
          v-if="vtuber.Twitch"
        >
          <font-awesome-icon
            :icon="['fab', 'twitch']"
            class="fa-fw"
            v-if="vtuber.Twitch"
          />
        </a>
        <a
          :href="`https://space.bilibili.com/${vtuber.BiliBili.SpaceID}`"
          target="_blank"
          class="vtuber-link__link hover:!text-bilibili"
          rel="noopener noreferrer"
          v-if="vtuber.BiliBili"
        >
          <font-awesome-icon
            :icon="['fab', 'bilibili']"
            class="fa-fw"
            v-if="vtuber.BiliBili"
          />
        </a>
        <a
          :href="`https://twitter.com/${vtuber.Twitter.Username}`"
          target="_blank"
          class="vtuber-link__link hover:!text-twitter"
          rel="noopener noreferrer"
          v-if="vtuber.Twitter"
        >
          <font-awesome-icon
            :icon="['fab', 'twitter']"
            class="fa-fw"
            v-if="vtuber.Twitter"
          />
        </a>
      </div>
    </div>
    <div class="card-vtuber-name">
      <router-link :to="`/vtuber/${vtuber.ID}`" class="card-vtuber-name__link">
        <h4 class="card-vtuber-name__title">
          {{ vtuber.EnName }}
        </h4>
        <span class="card-vtuber-name__nickname">
          {{ vtuber.NickName.toLowerCase() }}
        </span>
      </router-link>
    </div>
  </div>
</template>

<script>
// Add font awesome brands youtube, twitch, bilibili, twitter
import { library } from "@fortawesome/fontawesome-svg-core"
import { faSkull } from "@fortawesome/free-solid-svg-icons"
import {
  faYoutube,
  faTwitch,
  faBilibili,
  faTwitter,
} from "@fortawesome/free-brands-svg-icons"

library.add(faYoutube, faTwitch, faBilibili, faTwitter, faSkull)

export default {
  props: {
    vtuber: {
      type: Object,
    },
  },
  async created() {},
  computed: {
    liveLink() {
      if (this.vtuber.IsLive.Youtube) return this.vtuber.IsLive.Youtube.URL
      else if (this.vtuber.IsLive.Twitch) return this.vtuber.IsLive.Twitch.URL
      else return this.vtuber.IsLive.BiliBili.URL
    },
    randomAvatar() {},
    GroupName() {
      return (
        this.vtuber.Group.GroupName.charAt(0).toUpperCase() +
        this.vtuber.Group.GroupName.slice(1).replace("_", " ")
      )
    },
  },
}
</script>

<style lang="scss" scoped>
.card-vtuber {
  @apply bg-white dark:bg-slate-800 rounded-md overflow-hidden shadow-sm hover:bg-slate-100 dark:hover:bg-slate-900 hover:shadow-md dark:shadow-white/5 hover:scale-105 select-none transition duration-300 ease-in-out;

  &-image {
    @apply w-full aspect-square bg-smolame bg-cover relative;
  }

  &-name {
    &__link {
      @apply flex flex-col justify-center p-2;
    }

    &__title {
      @apply font-bold text-lg tracking-tight leading-6 truncate dark:text-white;
    }

    &__nickname {
      @apply text-xs text-stone-600 dark:text-stone-400;
    }
  }
}

.tag-vtuber {
  @apply absolute flex items-center bg-slate-100/80 dark:bg-slate-500/80 rounded-br-md text-xs select-none;

  &-agency {
    @apply px-[0.325rem] flex items-center space-x-2 justify-center h-6 cursor-pointer;

    &__icon {
      @apply w-5 object-contain drop-shadow-md rounded-md;
      image-rendering: smooth;
    }

    &__flag {
      @apply h-3 object-contain rounded-sm drop-shadow-md;
      image-rendering: smooth;
    }

    &__hover {
      @apply absolute top-6 px-2 py-1 -left-2 max-w-[10rem] truncate bg-slate-100/80 dark:bg-slate-500/80 font-semibold rounded-r-md scale-0 origin-top-left transition-transform ease-in-out duration-200;
    }

    &:hover &__hover {
      @apply scale-100;
    }
  }

  & > span {
    @apply hidden;
  }

  &-live {
    @apply bg-red-500 px-[0.325rem] h-6 rounded-br-md flex items-center space-x-1 font-semibold text-white;
  }
}

.vtuber-link {
  @apply absolute bottom-0 right-0 bg-slate-100/80  dark:bg-slate-500/80 rounded-tl-md px-[0.325rem] space-x-1;

  &__link {
    @apply inline-block py-[0.125rem] px-1 text-stone-700 dark:text-slate-50 transition-colors duration-200 ease-in-out;
  }
}

.inactive {
  .tag-vtuber {
    @apply z-[2];
  }

  .card-vtuber-image {
    @apply bg-rip bg-contain bg-no-repeat bg-gray-600;
    &__img {
      @apply grayscale opacity-40;
    }
  }
}
</style>
