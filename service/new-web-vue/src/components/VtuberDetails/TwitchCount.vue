<template>
  <div class="platform">
    <h4 class="platform-title">
      <font-awesome-icon
        :icon="['fab', 'twitch']"
        class="platform-title__icon fa-fw"
      />
      Twitch
    </h4>
    <a
      :href="`https://socialblade.com/twitch/user/${twitch.Username}/realtime`"
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
          <CountTo :endVal="twitch.Followers" />
        </span>
        <span class="card-title">Followers</span>
      </div>
      <div class="platform-card">
        <span class="card-count">
          <CountTo :endVal="twitch.ViwersCount" />
        </span>
        <span class="card-title">Viewers</span>
      </div>
    </div>
  </div>
  <hr class="m-2" />
</template>

<script>
import { CountTo } from "vue3-count-to"
import { library } from "@fortawesome/fontawesome-svg-core"
import { faTwitch } from "@fortawesome/free-brands-svg-icons"
import { faCircleRight } from "@fortawesome/free-solid-svg-icons"

library.add(faTwitch, faCircleRight)

export default {
  props: {
    twitch: {
      type: Object,
    },
  },
  components: {
    CountTo,
  },
}
</script>

<style lang="scss" scoped>
.platform {
  @apply mx-4 grid items-center;
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
    @apply text-lg font-semibold underline decoration-twitch decoration-2 underline-offset-4;
    grid-area: platform-title;
  }

  &-link {
    @apply justify-self-end py-2 text-right text-sm font-semibold hover:text-slate-700 hover:underline dark:hover:text-slate-300;
    grid-area: platform-link;
  }

  &-cards {
    @apply grid p-2;
    grid-template-columns: 1fr;
    grid-area: platform-card;

    @media (min-width: 640px) {
      grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    }
  }

  &-card {
    @apply m-1 inline-flex flex-col items-center justify-center rounded-md bg-twitch py-6 px-2;
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
