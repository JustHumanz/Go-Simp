<template>
  <header class="header">
    <div
      class="header-banner bg-cyan-300"
      :style="`background-image: url(${vtuber.Youtube.Banner})`"
      referrerpolicy="no-referrer"
      v-if="vtuber.Youtube"
    ></div>
    <div
      class="header-banner bg-cyan-300"
      :style="`background-image: url(${vtuber.BiliBili.Banner})`"
      referrerpolicy="no-referrer"
      v-else-if="vtuber.BiliBili"
    ></div>
    <div class="header-banner bg-cyan-300" v-else></div>
    <div class="header-info">
      <div class="header-profile-pic">
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
      </div>
      <div class="header-vtuber-name">
        <h4 class="header-vtuber-name__name">
          {{ vtuber.EnName }}{{ vtuber.JpName ? ` (${vtuber.JpName})` : "" }}
          <span class="nickname"> {{ vtuber.NickName }}</span>
        </h4>
        <div class="header-vtuber-name__group">
          <img
            :src="vtuber.Group.IconURL"
            :alt="vtuber.Group.GroupName"
            class="header-vtuber-name__group-icon"
            v-if="vtuber.Group.ID !== 10"
          />
          {{ vtuber.Group.ID === 10 ? "Vtuber" : vtuber.Group.GroupName }}
          {{ vtuber.Regions.name }}
        </div>
      </div>
    </div>
  </header>
</template>

<script>
import { library } from "@fortawesome/fontawesome-svg-core"

export default {
  props: {
    vtuber: Object,
  },
}
</script>

<style lang="scss" scoped>
.header {
  &-banner {
    @apply w-full h-[11.25rem] bg-center;
  }
  &-info {
    @apply flex items-center flex-col sm:flex-row;
  }
  &-profile-pic {
    @apply absolute -mt-[4.5rem] sm:mt-0 sm:ml-10 w-32 rounded-md shadow-md overflow-hidden;
  }
  &-vtuber-name {
    @apply w-full flex items-center sm:items-stretch mt-14 px-3 sm:mt-0 sm:ml-[11.5rem] flex-col text-center sm:text-left;

    &__name {
      @apply text-2xl font-semibold pt-3;
      .nickname {
        @apply text-gray-600 text-xl font-thin;
      }
    }

    &__group {
      @apply text-sm font-light text-gray-500;

      &-icon {
        @apply w-5 object-contain rounded-md inline-block;
      }
    }
  }
}
</style>
