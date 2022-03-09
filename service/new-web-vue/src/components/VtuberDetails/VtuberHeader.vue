<template>
  <header class="header">
    <div
      class="header-banner bg-cyan-300"
      :style="`background-image: url(${vtuber.Youtube.Banner.replace(
        's1200',
        ''
      )}s1707-fcrop64=1,00005a57ffffa5a8-k-c0xffffffff-no-nd-rj)`"
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
      <div
        class="header-profile-pic"
        :class="{ inactive: vtuber.Status == 'Inactive' }"
      >
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
        <span
          class="live-tag"
          v-if="
            vtuber.IsLive.Youtube ||
            vtuber.IsLive.Twitch ||
            vtuber.IsLive.BiliBili
          "
        >
          LIVE
        </span>
      </div>
      <div class="header-vtuber-name">
        <h4 class="header-vtuber-name__name">
          {{ vtuber.EnName }}{{ vtuber.JpName ? ` (${vtuber.JpName})` : "" }}
          <div class="nickname">
            {{ vtuber.NickName }}
            <span class="nickname-hover"
              >This nickname is used for initials when calling in the Discord
              command
              <router-link to="/docs#members">Learn More</router-link></span
            >
          </div>
        </h4>
        <div class="header-vtuber-name__group">
          <router-link :to="`/vtubers/${vtuber.Group.ID}`">
            <img
              :src="vtuber.Group.IconURL"
              :alt="vtuber.Group.GroupName"
              class="header-vtuber-name__group-icon"
              v-if="vtuber.Group.ID !== 10"
            />
            <img
              :src="`/src/assets/flags/${vtuber.Regions.flagCode}.svg`"
              :alt="vtuber.Group.GroupName"
              class="header-vtuber-name__group-icon"
              v-else
            />
            {{
              vtuber.Group.ID === 10
                ? "Vtuber"
                : vtuber.Group.GroupName.replace("_", " ")
            }}
            {{ vtuber.Regions.name }} </router-link
          ><router-link to="/vtubers?inac=true">
            {{ vtuber.Status === "Inactive" ? " (Inactive)" : "" }}
          </router-link>
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
.live-tag {
  @apply w-full inline-block bg-red-600 text-white font-semibold text-center text-xs py-1 absolute bottom-0 select-none;
}

.inactive {
  @apply bg-rip bg-contain bg-no-repeat bg-gray-600;
  .header-profile-pic__img {
    @apply grayscale opacity-40;
  }
}

.header {
  &-banner {
    @apply w-full /*&h-[11.25rem]*/ bg-center bg-cover bg-no-repeat;
    height: calc(16.1290322581vw - 1px);

    // @media (min-width: 640px) {
    //   height: calc(16.1290322581vw - 1px);
    // }

    @media (min-width: 768px) {
      height: calc((100vw - 240px) / 6.2 - 1px);
    }
  }
  &-info {
    @apply flex items-center flex-col sm:flex-row pb-2 relative;
  }
  &-profile-pic {
    @apply absolute -top-9 sm:-top-14 sm:ml-10 w-24 sm:w-32 rounded-md shadow-md overflow-hidden;
  }
  &-vtuber-name {
    @apply w-full flex items-center sm:items-stretch mt-14 px-3 sm:mt-0 sm:ml-[10.5rem] flex-col text-center sm:text-left;

    &__name {
      @apply text-xl font-semibold pt-3;
      .nickname {
        @apply text-gray-600 text-xl font-thin block sm:inline-block relative cursor-pointer;

        a {
          @apply text-white hover:text-gray-300;
        }

        &-hover {
          @apply bg-sky-300 text-sm font-normal px-2 py-1 rounded-md w-44 text-left inline-block absolute top-7 left-1/2 -translate-x-1/2 -translate-y-5 invisible opacity-0 delay-100 duration-300 ease-in-out;
          transition-property: opacity, transform;

          &::before {
            // add arrow up in center using tailwind
            @apply content-[''] border-x-8 border-b-8 border-solid border-x-transparent border-b-sky-300 absolute -top-1 left-1/2 -translate-x-1/2;
          }
        }

        &:hover {
          .nickname-hover {
            @apply opacity-100 translate-y-0 visible;
          }
        }
      }
    }

    &__group {
      @apply text-sm font-light text-gray-500 inline-block mt-px;

      a {
        @apply hover:text-gray-700;
      }

      &-icon {
        @apply w-5 h-5 object-contain rounded-md inline-block -mt-1;
      }
    }
  }
}
</style>
