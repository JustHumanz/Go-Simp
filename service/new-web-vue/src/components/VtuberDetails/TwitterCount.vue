<template>
  <div class="platform">
    <h4 class="platform-title">
      <font-awesome-icon
        :icon="['fab', 'twitter']"
        class="platform-title__icon fa-fw"
      />
      Twitter
    </h4>
    <a
      :href="`https://socialblade.com/twitter/user/${twitter.Username}/realtime`"
      class="platform-link"
      target="_blank"
    >
      Realtime Followers
      <font-awesome-icon
        icon="circle-right"
        class="platform-link__icon fa-fw"
      />
    </a>
    <div class="platform-cards">
      <div class="platform-card">
        <span class="card-count">
          <CountTo :endVal="twitter.Followers" />
        </span>
        <span class="card-title">Followers</span>
      </div>
      <a
        :href="`https://twitter.com/hashtag/${twitter.Fanart.replace('#', '')}`"
        target="_blank"
        class="platform-card"
        v-if="twitter.Fanart && twitter.Fanart !== ''"
      >
        <span class="card-count"> #{{ twitter.Fanart.replace("#", "") }} </span>
        <span class="card-title">Fanart Hastag</span>
      </a>
      <a
        :href="`https://twitter.com/hashtag/${twitter.LewdFanart.replace(
          '#',
          ''
        )}`"
        target="_blank"
        class="platform-card"
        v-if="twitter.LewdFanart && twitter.LewdFanart !== ''"
      >
        <span class="card-count">
          #{{ twitter.LewdFanart.replace("#", "") }}
        </span>
        <span class="card-title">NSFW Hastag</span>
      </a>
    </div>
  </div>
  <hr class="m-2" />
</template>

<script>
import { CountTo } from "vue3-count-to"
import { library } from "@fortawesome/fontawesome-svg-core"
import { faTwitter } from "@fortawesome/free-brands-svg-icons"

library.add(faTwitter)

export default {
  props: {
    twitter: {
      type: Object,
    },
  },
  components: {
    CountTo,
  },
}
</script>

<style lang="scss" scoped>
a.platform-card {
  @apply hover:brightness-90;
}
.platform {
  @apply grid mx-4 items-center;
  grid-template-columns: 1fr;
  grid-template-areas:
    "platform-title"
    "platform-card"
    "platform-link";

  @media (min-width: 640px) {
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    grid-template-areas:
      "platform-title platform-title platform-link"
      "platform-card platform-card platform-card";
  }

  &-title {
    @apply text-lg font-semibold underline decoration-2 underline-offset-4 decoration-twitter;
    grid-area: platform-title;
  }

  &-link {
    @apply text-sm font-semibold py-2 justify-self-end text-right hover:underline hover:text-slate-700 dark:hover:text-slate-300;
    grid-area: platform-link;
  }

  &-cards {
    @apply p-2 grid;
    grid-template-columns: 1fr;
    grid-area: platform-card;

    @media (min-width: 640px) {
      grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    }
  }

  &-card {
    @apply inline-flex flex-col items-center justify-center bg-twitter py-6 px-2 m-1 rounded-md;
  }
}

.card {
  &-count {
    @apply text-2xl font-semibold text-white;
  }

  &-title {
    @apply text-xs font-light text-gray-200;
  }
}
</style>
