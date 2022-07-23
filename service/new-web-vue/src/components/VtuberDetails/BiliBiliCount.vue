<template>
  <div class="platform">
    <h4 class="platform-title">
      <font-awesome-icon
        :icon="['fab', 'bilibili']"
        class="platform-title__icon fa-fw"
      />
      BiliBili
    </h4>
    <div class="platform-cards">
      <div class="platform-card">
        <span class="card-count">
          <CountTo :endVal="bilibili.Followers" />
        </span>
        <span class="card-title">Followers</span>
      </div>
      <div class="platform-card">
        <span class="card-count">
          <CountTo :endVal="bilibili.ViwersCount" />
        </span>
        <span class="card-title">Viewers</span>
      </div>
      <!-- <div class="platform-card">
        <span class="card-count">
          <CountTo :endVal="bilibili.TotalVideos" />
        </span>
        <span class="card-title">Videos/Archives</span>
      </div> -->
    </div>
  </div>
  <hr class="m-2" />
</template>

<script>
import { CountTo } from "vue3-count-to"
import { library } from "@fortawesome/fontawesome-svg-core"
import { faBilibili } from "@fortawesome/free-brands-svg-icons"
import { faCircleRight } from "@fortawesome/free-solid-svg-icons"

library.add(faBilibili, faCircleRight)

export default {
  props: {
    bilibili: {
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
    "platform-card";

  &-title {
    @apply text-lg font-semibold underline decoration-bilibili decoration-2 underline-offset-4;
    grid-area: platform-title;
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
    @apply m-1 inline-flex flex-col items-center justify-center rounded-md bg-bilibili py-6 px-2;
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
