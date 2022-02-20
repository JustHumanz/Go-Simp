<template>
  <div class="title">
    <span class="title__span">
      <font-awesome-icon :icon="['fas', 'circle-info']" class="title__svg" />
      Info Vtuber
    </span>
  </div>
  <section class="vtuber-view" v-if="member">
    <div class="profile">
      <div class="profile-photo">
        <img
          v-if="member.Youtube"
          referrerpolicy="no-referrer"
          :src="member.Youtube.Avatar"
          onerror="this.src='/src/assets/smolame.jpg'"
          :alt="member.EnName"
          class="photo-profile__img"
        />
        <img
          v-else-if="member.BiliBili"
          referrerpolicy="no-referrer"
          :src="member.BiliBili.Avatar"
          onerror="this.src='/src/assets/smolame.jpg'"
          :alt="member.EnName"
          class="photo-profile__img"
        />
        <img
          v-else-if="member.Twitch"
          referrerpolicy="no-referrer"
          :src="member.Twitch.Avatar"
          onerror="this.src='/src/assets/smolame.jpg'"
          :alt="member.EnName"
          class="photo-profile__img"
        />
        <img
          v-else
          referrerpolicy="no-referrer"
          src="/src/assets/smolame.jpg"
          :alt="member.EnName"
          class="photo-profile__img"
        />
      </div>
      <div class="profile-info">
        <h3 class="profile-name">
          {{ member.EnName }}{{ member.JpName ? `/${member.JpName}` : "" }}
          <span class="profile-name__span"> ({{ member.NickName }}) </span>
        </h3>

        <div class="profile-group">
          <img
            v-if="member.Group.ID !== 10"
            :src="member.Group.GroupIcon"
            :alt="member.Group.GroupName"
            class="profile-group__img"
          />
          <span>{{
            member.Group.ID !== 10
              ? member.Group.GroupName
              : "Vtuber Independent"
          }}</span>
        </div>

        <div class="profile-region">
          <img
            v-if="member.Regions.flagCode"
            :src="`/src/assets/flags/${member.Regions.flagCode}.svg`"
            :alt="member.Regions.name"
            class="profile-region__img"
          />
          <span>{{ member.Regions.name }}</span>
        </div>
        <div
          class="profile-live"
          v-if="member.IsBiliLive || member.IsYtLive || member.IsTwitchLive"
        >
          Watch live on
          <a :href="member.LiveURL" target="_blank" rel="noopener noreferrer">
            <!-- Add or when multisteam -->
            <span v-if="member.IsYtLive">
              <font-awesome-icon
                :icon="['fab', 'youtube']"
                class="profile-live__svg"
              />
              YouTube
            </span>
            <span
              v-if="
                (member.IsYtLive && member.IsBiliLive) ||
                (member.IsYtLive && member.IsTwitchLive) ||
                (member.IsBiliLive && member.IsYtLive && member.IsTwitchLive)
              "
            >
              or
            </span>
            <span v-if="member.IsBiliLive">
              <font-awesome-icon
                :icon="['fab', 'bilibili']"
                class="profile-live__svg"
              />
              Bilibili
            </span>
            <span
              v-if="
                (member.IsBiliLive && member.IsTwitchLive) ||
                (member.IsTwitchLive && member.IsYtLive && member.IsBiliLive)
              "
            >
              or
            </span>
            <span v-if="member.IsTwitchLive">
              <font-awesome-icon
                :icon="['fab', 'twitch']"
                class="profile-live__svg"
              />
              Twitch
            </span>
          </a>
        </div>
      </div>
    </div>
    <h4 class="sub-title">
      <font-awesome-icon
        :icon="['fas', 'chart-line']"
        class="sub-title__svg fa-fw"
      />
      Statistic
    </h4>
    <div class="followers">
      <div class="card followers-youtube" v-if="member.Youtube">
        <a
          class="card-title"
          :href="`https://youtube.com/channel/${member.Youtube.ID}`"
          target="_blank"
        >
          <font-awesome-icon
            :icon="['fab', 'youtube']"
            class="fa-fw card-title__svg"
          />
          YouTube
        </a>
        <div class="card-count">
          <CountTo :endVal="member.Youtube.Subscriber"></CountTo>
          <span class="card-count__span">Subscribers</span>
        </div>
        <a
          :href="`https://socialcounts.org/youtube-live-subscriber-count/${member.Youtube.ID}`"
          target="_blank"
          class="card-link"
          >Realtime</a
        >
      </div>
      <div class="card followers-bilibili" v-if="member.BiliBili">
        <a
          class="card-title"
          :href="`https://space.bilibili.com/${member.BiliBili.ID}`"
          target="_blank"
        >
          <font-awesome-icon
            :icon="['fab', 'bilibili']"
            class="fa-fw card-title__svg"
          />
          Bilibili
        </a>
        <div class="card-count">
          <CountTo :endVal="member.BiliBili.Follower"></CountTo>
          <span class="card-count__span">Followers</span>
        </div>
      </div>
      <div class="card followers-twitch" v-if="member.Twitch">
        <a
          class="card-title"
          :href="`https://www.twitch.tv/${member.Twitch.UserName}`"
          target="_blank"
        >
          <font-awesome-icon
            :icon="['fab', 'twitch']"
            class="fa-fw card-title__svg"
          />
          Twitch
        </a>
        <div class="card-count">
          <CountTo :endVal="member.Twitch.Followers"></CountTo>
          <span class="card-count__span">Followers</span>
        </div>
        <a
          :href="`https://livecounts.org/twitch-follower-count/${member.Twitch.UserName}`"
          target="_blank"
          class="card-link"
          >Realtime</a
        >
      </div>
      <div class="card followers-yt-views" v-if="member.Youtube">
        <a
          class="card-title"
          :href="`https://www.youtube.com/channel/${member.Youtube.ID}`"
          target="_blank"
        >
          <font-awesome-icon
            :icon="['fab', 'youtube']"
            class="fa-fw card-title__svg"
          />
          YouTube Views
        </a>
        <div class="card-count">
          <CountTo :endVal="member.Youtube.ViwersCount"></CountTo>
          <span class="card-count__span">Views</span>
        </div>
      </div>
      <div class="card followers-bili-views" v-if="member.BiliBili">
        <a
          class="card-title"
          :href="`https://space.bilibili.com/${member.BiliBili.ID}`"
          target="_blank"
        >
          <font-awesome-icon
            :icon="['fab', 'bilibili']"
            class="fa-fw card-title__svg"
          />
          Bilibili Views
        </a>
        <div class="card-count">
          <CountTo :endVal="member.BiliBili.ViwersCount"></CountTo>
          <span class="card-count__span">Views</span>
        </div>
      </div>
      <div class="card followers-twitch-views" v-if="member.Twitch">
        <a
          class="card-title"
          :href="`https://www.twitch.tv/${member.Twitch.UserName}`"
          target="_blank"
        >
          <font-awesome-icon
            :icon="['fab', 'twitch']"
            class="fa-fw card-title__svg"
          />
          Twitch Views
        </a>
        <div class="card-count">
          <CountTo :endVal="member.Twitch.ViwersCount"></CountTo>
          <span class="card-count__span">Views</span>
        </div>
      </div>
      <div class="card followers-twitter" v-if="member.Twitter">
        <a
          class="card-title"
          :href="`https://twitter.com/${member.Twitter.UserName}`"
          target="_blank"
        >
          <font-awesome-icon
            :icon="['fab', 'twitter']"
            class="fa-fw card-title__svg"
          />
          Twitter
        </a>
        <div class="card-count">
          <CountTo :endVal="member.Twitter.Followers"></CountTo>
          <span class="card-count__span">Followers</span>
        </div>
      </div>
      <div class="card followers-fans" v-if="member.Fanbase">
        <div class="card-title">
          <font-awesome-icon
            :icon="['fas', 'users']"
            class="fa-fw card-title__svg"
          />
          Fanbase Name
        </div>
        <span class="card-span">{{ member.Fanbase }}</span>
      </div>
    </div>
    <h4 class="sub-title">
      <font-awesome-icon
        :icon="['fas', 'paintbrush']"
        class="sub-title__svg fa-fw"
      />
      Fan Art
    </h4>
    <div class="fanart">
      <div class="card fanart-twitter" v-if="member.Twitter">
        <div class="card-title">
          <font-awesome-icon
            :icon="['fab', 'twitter']"
            class="fa-fw card-title__svg"
          />
          Fanart Hastag
        </div>
        <span class="card-span">{{ member.Twitter.Fanart }}</span>
        <a
          :href="`https://twitter.com/hashtag/${member.Twitter.Fanart.replace(
            '#',
            ''
          )}`"
          target="_blank"
          class="card-link"
          >Find</a
        >
      </div>
      <div class="card fanart-pixiv">
        <div class="card-title">
          <font-awesome-icon
            :icon="['fas', 'paintbrush']"
            class="fa-fw card-title__svg"
          />
          Fanart Pixiv
        </div>
        <span class="card-span">#{{ member.JpName || member.EnName }}</span>
        <a
          :href="`https://www.pixiv.net/en/tags/${
            member.JpName || member.EnName
          }/artworks`"
          target="_blank"
          class="card-link"
          >Find</a
        >
      </div>
    </div>
  </section>
  <section class="loading" v-if="!member">
    <img :src="`/src/assets/loading/${Math.floor(Math.random() * 7)}.gif`" class="loading__img" alt="ame loading">
  </section>
</template>

<style lang="scss" scoped>
.title {
  @apply text-2xl font-semibold uppercase bg-blue-400 py-3 w-full flex flex-wrap;

  &__span {
    @apply w-[90%] md:w-[70%] lg:w-[65%] mx-auto text-white;
  }
  &__svg {
    @apply text-blue-200;
  }
}

.sub-title {
  @apply text-xl font-semibold uppercase py-3 w-full;
}
.vtuber-view {
  @apply w-full md:w-[80%] lg:w-[75%] mx-auto px-2;

  .profile {
    @apply flex flex-col sm:flex-row items-center sm:items-stretch py-4;

    &-photo {
      @apply object-contain h-32 aspect-square bg-smolame bg-cover overflow-hidden rounded-lg;
    }

    &-info {
      @apply px-3 py-2 sm:py-0 grid justify-items-center text-center gap-y-2;
      grid-template-columns: 1fr 1fr;
      grid-template-areas:
        "name name"
        "group region"
        "live live";

      // media query small
      @media (min-width: 640px) {
        @apply justify-items-stretch text-left content-start gap-y-1;
        justify-items: auto;
        grid-template-columns: 1fr;
        grid-template-areas:
          "name"
          "group"
          "region"
          "live";
      }
    }

    &-name {
      @apply text-xl font-semibold;
      grid-area: name;

      &__span {
        @apply font-thin text-gray-600;
      }
    }

    &-group {
      grid-area: group;
    }

    &-region {
      grid-area: region;
    }

    &-live {
      @apply bg-red-500 text-white justify-self-stretch text-center text-sm rounded-md p-1;
      grid-area: live;
    }

    &-group,
    &-region {
      @apply flex h-6 text-xs items-center;

      &__img {
        @apply h-6 w-6 object-contain mr-2 drop-shadow-md;
      }
    }
  }
}

.followers {
  @apply grid gap-3 pb-4;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));

  &-youtube,
  &-yt-views {
    @apply bg-youtube;
  }

  &-bilibili,
  &-bili-views {
    @apply bg-bilibili;
  }

  &-twitch,
  &-twitch-views {
    @apply bg-twitch;
  }

  &-twitter {
    @apply bg-twitter;
  }

  &-fans {
    @apply bg-blue-400;
  }
}

.fanart {
  @apply grid gap-3 pb-5 lg:pb-7;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));

  &-twitter {
    @apply bg-twitter;
  }

  &-pixiv {
    @apply bg-pixiv;
  }
}

.card {
  @apply text-center rounded-lg overflow-hidden shadow-md border-2 border-black/30 relative flex justify-center items-center flex-col h-44;

  &-title {
    @apply text-base font-semibold bg-black/25 text-white block w-full py-2 absolute top-0;

    &__svg {
      @apply text-slate-300/80;
    }
  }

  &-count {
    @apply text-2xl font-semibold text-white inline-block w-full leading-3 pt-5 pb-3;

    &__span {
      @apply text-gray-100 text-sm block mt-2;
    }

    // add more padding when no have.card-link in bottom
    &:last-child {
      @apply pt-8;
    }
  }

  &-span {
    @apply text-white text-2xl font-semibold block;
  }

  &-link {
    @apply absolute bottom-0 text-sm font-semibold text-gray-100 bg-white/25 block w-full py-2;
  }
}

.loading {
  @apply flex justify-center items-center h-screen;

  &__img {
    @apply object-contain h-32 aspect-square;
  }
}
</style>

<script>
import axios from "axios"
import Config from "../config.json"
import regionConfig from "../region.json"
import { CountTo } from "vue3-count-to"
import { library } from "@fortawesome/fontawesome-svg-core"

import {
  faYoutube,
  faBilibili,
  faTwitch,
  faTwitter,
} from "@fortawesome/free-brands-svg-icons"
import {
  faCircleInfo,
  faPaintbrush,
  faUsers,
  faChartLine,
} from "@fortawesome/free-solid-svg-icons"

library.add(
  faYoutube,
  faBilibili,
  faTwitch,
  faTwitter,
  faCircleInfo,
  faPaintbrush,
  faUsers,
  faChartLine
)

export default {
  data() {
    return {
      member: null,
      fanart: null,
    }
  },
  components: {
    CountTo,
  },
  async created() {
    document.title = "Vtuber Details"

    const member_data = await axios
      .get(Config.REST_API + "/members/" + this.$route.params.id)
      .then((response) => response.data[0])

    member_data.Group = await axios
      .get(Config.REST_API + "/groups/" + member_data.GroupID)
      .then((response) => response.data[0])

    // delete GroupName and GroupID
    delete member_data.GroupName
    delete member_data.GroupID

    regionConfig.forEach((region) => {
      if (region.code === member_data.Region) {
        member_data.Regions = region
      }
    })

    this.member = member_data
    console.log(this.member)

    // get fanart
    this.fanart = await axios
      .get(Config.REST_API + "/fanart/random/member/" + this.$route.params.id)
      .then((response) => response.data)
    console.log(this.fanart)

    document.title = this.member.EnName + " - Vtuber Details"
  },
}
</script>
