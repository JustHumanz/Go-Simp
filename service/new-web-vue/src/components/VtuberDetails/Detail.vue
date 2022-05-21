<script setup>
import PixivLogo from "../icons/PixivLogo.vue"
</script>

<template>
  <div class="detail">
    <div class="detail-link">
      <a
        :href="`https://youtube.com/channel/${vtuber.Youtube.YoutubeID}`"
        v-if="vtuber.Youtube"
        target="_blank"
        rel="noopener noreferrer"
        class="detail-link__link bg-youtube"
      >
        <font-awesome-icon
          :icon="['fab', 'youtube']"
          size="lg"
          class="detail-link__link-icon"
        ></font-awesome-icon>
        <span class="detail-link__link-text">YouTube</span>
      </a>
      <a
        :href="`https://twitch.tv/${vtuber.Twitch.Username}`"
        v-if="vtuber.Twitch"
        target="_blank"
        rel="noopener noreferrer"
        class="detail-link__link bg-twitch"
      >
        <font-awesome-icon
          :icon="['fab', 'twitch']"
          size="lg"
          class="detail-link__link-icon"
        ></font-awesome-icon>
        <span class="detail-link__link-text">Twitch</span>
      </a>
      <a
        :href="`https://live.bilibili.com/${vtuber.BiliBili.LiveID}`"
        v-if="vtuber.BiliBili"
        target="_blank"
        rel="noopener noreferrer"
        class="detail-link__link bg-bilibili"
      >
        <font-awesome-icon
          :icon="['fab', 'bilibili']"
          size="lg"
          class="detail-link__link-icon"
        ></font-awesome-icon>
        <span class="detail-link__link-text">Live</span>
      </a>
      <a
        :href="`https://space.bilibili.com/${vtuber.BiliBili.SpaceID}`"
        v-if="vtuber.BiliBili"
        target="_blank"
        rel="noopener noreferrer"
        class="detail-link__link bg-bilibili"
      >
        <font-awesome-icon
          :icon="['fab', 'bilibili']"
          size="lg"
          class="detail-link__link-icon"
        ></font-awesome-icon>
        <span class="detail-link__link-text">Space</span>
      </a>
      <a
        class="detail-link__link bg-twitter"
        :href="`https://twitter.com/${vtuber.Twitter.Username}`"
        v-if="vtuber.Twitter"
        target="_blank"
        rel="noopener noreferrer"
      >
        <font-awesome-icon
          :icon="['fab', 'twitter']"
          size="lg"
          class="detail-link__link-icon"
        ></font-awesome-icon>
        <span class="detail-link__link-text"
          >@{{ vtuber.Twitter.Username }}</span
        >
      </a>
      <a
        :href="`https://twitter.com/hashtag/${vtuber.Twitter.Fanart.replace(
          '#',
          ''
        )}`"
        v-if="vtuber.Twitter && vtuber.Twitter.Fanart"
        target="_blank"
        rel="noopener noreferrer"
        class="detail-link__link bg-twitter"
      >
        <font-awesome-icon
          :icon="['fab', 'twitter']"
          size="lg"
          class="detail-link__link-icon"
        ></font-awesome-icon>
        <span class="detail-link__link-text"
          >Art (#{{ vtuber.Twitter.Fanart.replace("#", "") }})</span
        >
      </a>
      <a
        :href="`https://www.pixiv.net/en/tags/${fanart_pixiv}`"
        target="_blank"
        rel="noopener noreferrer"
        class="detail-link__link bg-pixiv"
      >
        <PixivLogo class="detail-link__link-icon" />
        <span class="detail-link__link-text"> Art (#{{ fanart_pixiv }}) </span>
      </a>
    </div>
    <div class="detail-hr"></div>
    <div class="detail-info">
      <div class="detail-info-item">
        <span class="detail-info-item__title"> Name </span>
        <span class="detail-info-item__value">
          {{ vtuber.EnName }}
        </span>
      </div>

      <div class="detail-info-item" v-if="vtuber.JpName">
        <span class="detail-info-item__title"> Japanese Name </span>
        <span class="detail-info-item__value">
          {{ vtuber.JpName }}
        </span>
      </div>

      <div class="detail-info-item">
        <span class="detail-info-item__title"> Group/Agency </span>
        <router-link
          :to="`/vtubers/${vtuber.Group.ID}`"
          class="detail-info-item__value"
        >
          <img
            :src="vtuber.Group.IconURL"
            :alt="GroupName"
            class="detail-info-item__icon"
            v-if="vtuber.Group.ID !== 10"
          />
          {{ GroupName }}
        </router-link>
      </div>

      <div class="detail-info-item">
        <span class="detail-info-item__title"> Region </span>
        <router-link
          :to="`/vtubers?reg=${vtuber.Region}`"
          class="detail-info-item__value"
        >
          <img
            :src="`/assets/flags/${vtuber.Regions.flagCode}.svg`"
            :alt="vtuber.Regions.name"
            class="detail-info-item__icon"
          />
          {{ vtuber.Regions.name }}
        </router-link>
      </div>

      <div class="detail-info-item">
        <span class="detail-info-item__title"> Status </span>
        <span class="detail-info-item__value">
          {{ vtuber.Status }}
        </span>
      </div>

      <div class="detail-info-item" v-if="vtuber.Fanbase">
        <span class="detail-info-item__title"> Fanbase </span>
        <span class="detail-info-item__value">
          {{ vtuber.Fanbase }}
        </span>
      </div>

      <div class="detail-info-item">
        <span class="detail-info-item__title"></span>
        <a href="#" onclick="return false" class="detail-info-item__value">
          <font-awesome-icon icon="triangle-exclamation" class="fa-fw mr-1" />
          Send issue
        </a>
      </div>
    </div>
  </div>
  <hr class="m-2" />
</template>

<script>
export default {
  props: {
    vtuber: Object,
  },
  computed: {
    GroupName() {
      return (
        this.vtuber.Group.GroupName.charAt(0).toUpperCase() +
        this.vtuber.Group.GroupName.slice(1).replace("_", " ")
      )
    },
    fanart_pixiv() {
      return this.vtuber.JpName
        ? this.vtuber.JpName.split("/")[0]
            .split(" ")
            .reduce((acc, cur) => `${acc}${acc !== "" ? "・" : ""}${cur}`, "")
        : this.vtuber.EnName.split(" ").reduce((acc, cur) => acc + cur, "")
    },
  },
}
</script>

<style lang="scss" scoped>
.detail {
  @apply mx-2 grid gap-2 transition-colors duration-100 ease-in-out md:mx-0;

  @media (min-width: 768px) {
    grid-template-columns: 1.75fr auto 1fr;
  }

  &-link {
    @apply flex h-max flex-wrap justify-center self-center text-sm font-semibold;

    &__link {
      @apply m-0.5 flex items-center justify-center rounded-full px-2.5 py-1 text-white transition duration-300 ease-in-out hover:brightness-90;

      &-icon {
        @apply mr-2;
      }

      &-text {
        @apply text-center;
      }
    }
  }

  &-info {
    @apply flex flex-col;

    &-item {
      @apply flex flex-row items-center justify-between;

      &__title {
        @apply text-sm font-light text-gray-600 dark:text-gray-300;
      }

      &__icon {
        @apply mr-2 w-[1em] drop-shadow-sm dark:shadow-white/10;
      }

      a {
        @apply hover:text-gray-600 dark:hover:text-gray-400;
      }

      &__value {
        @apply flex items-center text-sm text-gray-800 dark:text-gray-200;
      }
    }
  }

  &-hr {
    @apply border-t-[1px] border-gray-200 md:border-t-0 md:border-l-[1px];
  }

  div:not(.detail-hr) {
    @apply px-2;
  }
}
</style>
