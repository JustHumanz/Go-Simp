<template>
  <div class="card-vtuber" :class="{ inactive: vtuber.Status == 'Inactive' }">
    <div class="card-vtuber-image">
      <div class="tag-vtuber">
        <div class="tag-vtuber-agency">
          <img
            class="tag-vtuber-agency__icon"
            :src="vtuber.Group.IconURL"
            alt=""
            v-if="vtuber.Group.ID !== 10"
          />
          <img
            class="tag-vtuber-agency__flag"
            :src="`/src/assets/flags/${vtuber.Regions.flagCode}.svg`"
            :alt="vtuber.Regions.name"
          />
          <font-awesome-icon
            class="fa-fw"
            icon="skull"
            v-if="vtuber.Status == 'Inactive'"
          />
        </div>
        <div
          class="tag-vtuber-live"
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
          <span>LIVE</span>
        </div>
      </div>
      <router-link
        :to="`/vtuber/members/${vtuber.ID}`"
        class="card-vtuber-image__link"
      >
        <img
          class="card-vtuber-image__img"
          v-if="vtuber.Youtube"
          v-bind:src="vtuber.Youtube.Avatar.replace('s800', 's360')"
          referrerpolicy="no-referrer"
          onerror="this.src='/src/assets/smolame.jpg'"
          alt="Card image cap"
        />
        <img
          class="card-vtuber-image__img"
          v-else-if="vtuber.BiliBili"
          v-bind:src="`${vtuber.BiliBili.Avatar}@360w_360h_1c_1s.jpg`"
          referrerpolicy="no-referrer"
          onerror="this.src='/src/assets/smolame.jpg'"
          alt="Card image cap"
        />
        <img
          class="card-vtuber-image__img"
          v-else-if="vtuber.Twitch"
          v-bind:src="vtuber.Twitch.Avatar"
          referrerpolicy="no-referrer"
          onerror="this.src='/src/assets/smolame.jpg'"
          alt="Card image cap"
        />
        <img
          class="card-vtuber-image__img"
          v-else
          src="/src/assets/smolame.jpg"
          alt="Card image cap"
        />
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
  created() {
    // console.log("Show Vtuber: " + this.vtuber.EnName)
  },
}
</script>

<style lang="scss" scoped>
.card-vtuber {
  @apply bg-white rounded-md overflow-hidden shadow-sm hover:bg-slate-100 hover:shadow-md hover:scale-110 transition select-none;

  &-image {
    @apply w-full aspect-square bg-smolame bg-cover;
  }
}

.tag-vtuber {
  @apply absolute flex items-center bg-white/80 rounded-br-md overflow-hidden text-xs;

  &-agency {
    @apply px-[0.325rem] flex items-center space-x-1 justify-center h-6;

    &__icon {
      @apply w-5 object-contain drop-shadow-md;
    }

    &__flag {
      @apply h-3 object-contain rounded-sm drop-shadow-md;
    }
  }

  &-live {
    @apply bg-red-500 px-[0.325rem] h-6 rounded-br-md flex items-center space-x-1 font-semibold text-white;
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
