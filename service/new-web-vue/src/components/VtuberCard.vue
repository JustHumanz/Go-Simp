<template>
  <div class="container-vtuber">
    <div
      class="card-vtuber"
      :class="{ inactive: vtuber.Status == 'Inactive' }"
      :data-id="vtuber.ID"
    >
      <div class="card-vtuber-image">
        <div class="tag-vtuber">
          <div class="tag-vtuber-agency">
            <router-link
              :to="`/vtubers/` + vtuber.Group.ID"
              @click="getVtuberData(vtuber.Group.ID)"
            >
              <img
                draggable="false"
                class="tag-vtuber-agency__icon"
                :src="vtuber.Group.IconURL"
                :alt="vtuber.Group.GroupName"
                v-if="
                  $route.params.id != vtuber.Group.ID && vtuber.Group.ID !== 10
                "
              />
              <img
                draggable="false"
                class="tag-vtuber-agency__flag"
                :src="`/assets/flags/${vtuber.Regions?.code.toLowerCase()}.svg`"
                :alt="vtuber.Regions.name"
                onerror="this.src='/assets/flags/none.svg'"
                v-else-if="vtuber.Regions"
              />
              <img
                draggable="false"
                class="tag-vtuber-agency__flag"
                src="/assets/flags/none.svg"
                :alt="vtuber.Region"
                v-else
              />
            </router-link>

            <font-awesome-icon
              class="fa-fw"
              icon="skull"
              v-if="vtuber.Status == 'Inactive'"
            />
            <span class="tag-vtuber-agency__hover">{{
              `${agencyName} ${
                vtuber.Status === "Inactive" ? ` (Inactive)` : ""
              }`
            }}</span>
          </div>
          <div class="tag-vtuber-link">
            <a
              class="tag-vtuber-link__item hover:!text-youtube"
              :class="{ live: vtuber.IsLive.Youtube }"
              :href="ytLink"
              target="_blank"
              v-if="vtuber.Youtube"
            >
              <font-awesome-icon class="fa-fw" :icon="['fab', 'youtube']" />
            </a>
            <a
              class="tag-vtuber-link__item hover:!text-twitch"
              :class="{ live: vtuber.IsLive.Twitch }"
              :href="`https://twitch.tv/${vtuber.Twitch.Username}`"
              target="_blank"
              v-if="vtuber.Twitch"
            >
              <font-awesome-icon class="fa-fw" :icon="['fab', 'twitch']" />
            </a>
            <a
              class="tag-vtuber-link__item hover:!text-bilibili"
              :class="{ live: vtuber.IsLive.BiliBili }"
              :href="blLink"
              target="_blank"
              v-if="vtuber.BiliBili"
            >
              <font-awesome-icon class="fa-fw" :icon="['fab', 'bilibili']" />
            </a>
            <a
              class="tag-vtuber-link__item hover:!text-twitter"
              :href="`https://twitter.com/${vtuber.Twitter.Username}`"
              target="_blank"
              v-if="vtuber.Twitter"
            >
              <font-awesome-icon class="fa-fw" :icon="['fab', 'twitter']" />
            </a>
          </div>
        </div>
        <router-link
          :to="`/vtuber/${vtuber.ID}`"
          class="card-vtuber-image__link"
        >
          <img
            draggable="false"
            class="card-vtuber-image__img"
            :src="avatar"
            referrerpolicy="no-referrer"
            onerror="this.src='/assets/smolame.jpg'"
            :alt="vtuber.EnName"
          />
        </router-link>
        <div class="vtuber-link" v-if="type !== 'name' && type !== 'jpname'">
          <font-awesome-icon
            class="fa-fw"
            :icon="['fab', 'youtube']"
            v-if="type.includes('youtube')"
          />
          <font-awesome-icon
            class="fa-fw"
            :icon="['fab', 'twitch']"
            v-if="type.includes('twitch')"
          />
          <font-awesome-icon
            class="fa-fw"
            :icon="['fab', 'bilibili']"
            v-if="type.includes('bilibili')"
          />
          <font-awesome-icon
            class="fa-fw"
            :icon="['fab', 'twitter']"
            v-if="type.includes('twitter')"
          />
          <span>{{ `${countSort} ${textCount}` }}</span>
        </div>
      </div>
      <div class="card-vtuber-name">
        <router-link
          :to="`/vtuber/${vtuber.ID}`"
          class="card-vtuber-name__link"
        >
          <h4 class="card-vtuber-name__title">
            {{ titleName }}
          </h4>
          <span class="card-vtuber-name__nickname">
            {{ vtuber.NickName.toLowerCase() }}
          </span>
        </router-link>
      </div>
    </div>
    <div class="context-menu-vtuber">
      <div class="context-menu-vtuber__item" @click="copyToClipboard">
        <font-awesome-icon class="fa-fw" :icon="['fas', 'copy']" />
        <span>Show Youtube Channel</span>
      </div>
      <div class="context-menu-vtuber__item" @click="copyToClipboard">
        <font-awesome-icon class="fa-fw" :icon="['fas', 'copy']" />
        <span>Copy Nickname & Group</span>
      </div>
      <div class="context-menu-vtuber__item" @click="copyToClipboard">
        <font-awesome-icon class="fa-fw" :icon="['fas', 'copy']" />
        <span>Copy Group</span>
      </div>
      <div class="context-menu-vtuber__item" @click="copyToClipboard">
        <font-awesome-icon class="fa-fw" :icon="['fas', 'copy']" />
        <span>Copy</span>
      </div>
      <div class="context-menu-vtuber__item" @click="copyToClipboard">
        <font-awesome-icon class="fa-fw" :icon="['fas', 'copy']" />
        <span>Copy</span>
      </div>
      <div class="context-menu-vtuber__item" @click="copyToClipboard">
        <font-awesome-icon class="fa-fw" :icon="['fas', 'copy']" />
        <span>Copy</span>
      </div>
      <div class="context-menu-vtuber__item" @click="copyToClipboard">
        <font-awesome-icon class="fa-fw" icon="warning" />
        <span>Send Issues</span>
      </div>
    </div>
  </div>
</template>

<script>
// Add font awesome brands youtube, twitch, bilibili, twitter
import { library } from "@fortawesome/fontawesome-svg-core"
import { faSkull, faCopy, faWarning } from "@fortawesome/free-solid-svg-icons"
import {
  faYoutube,
  faTwitch,
  faBilibili,
  faTwitter,
} from "@fortawesome/free-brands-svg-icons"

library.add(
  faYoutube,
  faTwitch,
  faBilibili,
  faTwitter,
  faSkull,
  faCopy,
  faWarning
)

import { useMemberStore } from "@/stores/members.js"

export default {
  props: {
    vtuber: {
      type: Object,
    },
  },
  mounted() {
    const cardsContainer = document.querySelectorAll(".container-vtuber")

    cardsContainer.forEach((c) => {
      const card = c.children[0]
      const contextMenu = c.children[1]
    })
  },
  computed: {
    avatar() {
      return this.vtuber.Youtube
        ? this.vtuber.Youtube.Avatar.replace("s800", "s360")
        : this.vtuber.BiliBili
        ? `${this.vtuber.BiliBili.Avatar}@360w_360h_1c_1s.jpg`
        : this.vtuber.Twitch
        ? this.vtuber.Twitch.Avatar
        : this.vtuber.Twitter
        ? this.vtuber.Twitter.Avatar
        : "/assets/smolame.jpg"
    },
    GroupName() {
      return (
        this.vtuber.Group.GroupName.charAt(0).toUpperCase() +
        this.vtuber.Group.GroupName.slice(1).replace("_", " ")
      )
    },
    ytLink() {
      if (this.vtuber.IsLive?.Youtube) return this.vtuber.IsLive.Youtube.URL
      else if (this.vtuber.Youtube)
        return `https://youtube.com/channel/${this.vtuber.Youtube.YoutubeID}`
      else return null
    },
    blLink() {
      if (this.vtuber.IsLive.BiliBili) return this.vtuber.IsLive.BiliBili.URL
      else if (this.vtuber.BiliBili)
        return `https://space.bilibili.com/${this.vtuber.BiliBili.SpaceID}`
      else return null
    },
    type() {
      return useMemberStore().sorting.type
    },
    titleName() {
      if (useMemberStore().sorting.type === "jpname")
        return this.vtuber.JpName ? this.vtuber.JpName : this.vtuber.EnName
      else return this.vtuber.EnName
    },
    agencyName() {
      if (this.vtuber.Agency) {
        if (this.vtuber.Agency.match(/^\$/g))
          return this.vtuber.Agency.replace(/^\$/, "")
        else
          return `${this.vtuber.Agency} (${this.vtuber.Group.GroupName} ${this.vtuber.Region})`
      } else
        return `${
          this.vtuber.Group.ID === 10 ? `Vtuber` : this.vtuber.Group.GroupName
        } ${this.vtuber.Regions?.name || this.vtuber.Region}`
    },
    countSort() {
      const type = useMemberStore().sorting.type

      if (type === "youtube")
        return this.vtuber.Youtube
          ? this.FormatNumber(this.vtuber.Youtube.Subscriber)
          : 0
      else if (type === "bilibili")
        return this.vtuber.BiliBili
          ? this.FormatNumber(this.vtuber.BiliBili.Followers)
          : 0
      else if (type === "twitch")
        return this.vtuber.Twitch
          ? this.FormatNumber(this.vtuber.Twitch.Followers)
          : 0
      else if (type === "twitter")
        return this.vtuber.Twitter
          ? this.FormatNumber(this.vtuber.Twitter.Followers)
          : 0
      else if (type === "youtube_views")
        return this.vtuber.Youtube
          ? this.FormatNumber(this.vtuber.Youtube.ViwersCount)
          : 0
      else if (type === "bilibili_views")
        return this.vtuber.BiliBili
          ? this.FormatNumber(this.vtuber.BiliBili.ViwersCount)
          : 0
      else return 0
    },
    textCount() {
      const type = useMemberStore().sorting.type

      if (type === "youtube") return "Subscribers"
      else if (type.includes("views")) return "Views"
      else return "Followers"
    },
  },
  methods: {
    FormatNumber(num) {
      // Ensure number has max 3 significant digits (no rounding up can happen)

      if (num >= 1000000000) return (num / 1000000000).toFixed(2) + "B"
      else if (num >= 1000000) return (num / 1000000).toFixed(2) + "M"
      else if (num >= 1000) return (num / 1000).toFixed(2) + "K"
      else return num
    },
    async getVtuberData(id) {
      await useMemberStore().fetchMembers(id || null)
      useMemberStore().filterMembers()
      useMemberStore().sortingMembers()
    },
  },
}
</script>

<style lang="scss" scoped>
.container-vtuber {
  @apply relative;
}
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
  @apply absolute flex w-full select-none items-center justify-between text-sm;

  &-agency {
    @apply flex h-6 cursor-pointer items-center justify-center space-x-2 rounded-br-md bg-slate-100/75 px-[0.325rem] dark:bg-slate-500/80;

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

  &-link {
    @apply flex h-6 items-center overflow-hidden rounded-bl-md bg-slate-100/75 font-semibold text-slate-600 dark:bg-slate-500/80 dark:text-white;

    // check when no child
    &:empty {
      @apply px-0;
    }

    &__item {
      @apply h-full px-1 py-[3px] transition-colors duration-200 ease-in-out first:pl-[6px] last:pr-[6px] hover:bg-slate-600/20 dark:hover:bg-white/20;

      &.live {
        @apply bg-red-600 text-white hover:bg-red-300;
      }
    }
  }
}

.vtuber-link {
  @apply absolute bottom-0 right-0 space-x-1 rounded-tl-md bg-slate-100/80 px-[0.325rem] text-sm dark:bg-slate-500/80;
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

.context-menu-vtuber {
  @apply absolute top-0 z-[4] w-full rounded-md bg-slate-700 shadow-md;

  &__item {
    @apply flex cursor-pointer space-x-1 px-2 py-1.5 text-sm first:pt-2 hover:bg-slate-600/40;
  }
}
</style>
