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
            :src="`/assets/flags/${vtuber.Regions?.code.toLowerCase()}.svg`"
            :alt="vtuber.Regions.name"
            onerror="this.src='/assets/flags/none.svg'"
            v-if="vtuber.Regions"
          />
          <img
            draggable="false"
            class="tag-vtuber-agency__flag"
            src="/assets/flags/none.svg"
            :alt="vtuber.Region"
            v-else
          />

          <font-awesome-icon
            class="fa-fw"
            icon="skull"
            v-if="vtuber.Status == 'Inactive'"
          />
          <span class="tag-vtuber-agency__hover">{{
            `${vtuber.Group.ID === 10 ? `Vtuber` : GroupName} ${
              vtuber.Regions?.name || vtuber.Region
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
          onerror="this.src='/assets/smolame.jpg'"
          alt="Card image cap"
        />
        <img
          draggable="false"
          class="card-vtuber-image__img"
          v-else-if="vtuber.BiliBili"
          v-bind:src="`${vtuber.BiliBili.Avatar}@360w_360h_1c_1s.jpg`"
          referrerpolicy="no-referrer"
          onerror="this.src='/assets/smolame.jpg'"
          alt="Card image cap"
        />
        <img
          draggable="false"
          class="card-vtuber-image__img"
          v-else-if="vtuber.Twitch"
          v-bind:src="vtuber.Twitch.Avatar"
          referrerpolicy="no-referrer"
          onerror="this.src='/assets/smolame.jpg'"
          alt="Card image cap"
        />
        <img
          draggable="false"
          class="card-vtuber-image__img"
          v-else
          src="/assets/smolame.jpg"
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
  @apply select-none overflow-hidden rounded-md bg-white shadow-sm transition duration-300 ease-in-out hover:scale-105 hover:bg-slate-100 hover:shadow-md dark:bg-slate-500 dark:shadow-white/5 dark:hover:bg-slate-700;

  &-image {
    @apply relative aspect-square w-full bg-smolame bg-cover;
  }

  &-name {
    &__link {
      @apply flex flex-col justify-center p-2;
    }

    &__title {
      @apply truncate text-lg font-bold leading-6 tracking-tight dark:text-white;
    }

    &__nickname {
      @apply text-xs text-stone-600 dark:text-stone-300;
    }
  }
}

.tag-vtuber {
  @apply absolute flex select-none items-center rounded-br-md bg-slate-100/80 text-xs dark:bg-slate-500/80;

  &-agency {
    @apply flex h-6 cursor-pointer items-center justify-center space-x-2 px-[0.325rem];

    &__icon {
      @apply w-5 rounded-md object-contain drop-shadow-md;
      image-rendering: smooth;
    }

    &__flag {
      @apply h-3 rounded-sm object-contain drop-shadow-md;
      image-rendering: smooth;
    }

    &__hover {
      @apply absolute top-6 -left-2 max-w-[10rem] origin-top-left scale-0 truncate rounded-r-md bg-slate-100/80 px-2 py-1 font-semibold transition-transform duration-200 ease-in-out dark:bg-slate-500/80;
    }

    &:hover &__hover {
      @apply scale-100;
    }
  }

  & > span {
    @apply hidden;
  }

  &-live {
    @apply flex h-6 items-center space-x-1 rounded-br-md bg-red-500 px-[0.325rem] font-semibold text-white;
  }
}

.vtuber-link {
  @apply absolute bottom-0 right-0 space-x-1  rounded-tl-md bg-slate-100/80 px-[0.325rem] dark:bg-slate-500/80;

  &__link {
    @apply inline-block py-[0.125rem] px-1 text-stone-700 transition-colors duration-200 ease-in-out dark:text-slate-50;
  }
}

.inactive {
  .tag-vtuber {
    @apply z-[2];
  }

  .card-vtuber-image {
    @apply bg-gray-600 bg-rip bg-contain bg-no-repeat;
    &__img {
      @apply opacity-40 grayscale;
    }
  }
}
</style>
