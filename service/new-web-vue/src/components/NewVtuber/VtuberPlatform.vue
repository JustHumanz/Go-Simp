<template>
  <div class="vtuber-form show" ref="form" :data-id="id">
    <a
      class="vtuber-link"
      href="#"
      onclick="return false"
      @click="toggleContent"
    >
      <span class="vtuber-link__text">{{ getVtuberName }}</span>
      <span class="vtuber-link__icon">
        <a
          class="delete-vtuber"
          href="#"
          onclick="return false"
          @click="deletePlatform"
          ><font-awesome-icon icon="trash-can" class="fa-fw"
        /></a>
        <font-awesome-icon icon="chevron-down" class="arrow" />
      </span>
    </a>
    <div class="vtuber__content" ref="content">
      <div class="vtuber__content-item">
        <label for="name">Nickname</label>
        <input type="text" id="name" name="name" autocomplete="off" />
        <small class="description">
          A short name/calling in vtuber for command element</small
        >
        <small class="error"></small>
      </div>

      <div class="vtuber__content-item">
        <label for="en-name">Vtuber Name</label>
        <input
          type="text"
          id="en-name"
          name="en-name"
          autocomplete="off"
          v-model="vtuberName"
        />
        <small class="description"> Name Vtuber in English </small>
        <small class="error"></small>
      </div>

      <div class="vtuber__content-item">
        <label for="jp-name">Japanese Name</label>
        <input type="text" id="jp-name" name="jp-name" autocomplete="off" />
        <small class="description"
          >(Optional) Name Vtuber in Japanese, for easier in pixiv fanart</small
        >
        <small class="error"></small>
      </div>

      <div class="vtuber__content-item">
        <label for="fanbase">Fanbase Name</label>
        <input type="text" id="fanbase" name="fanbase" autocomplete="off" />
        <small class="description">(Optional) Nickname for vtuber fans</small>
        <small class="error"></small>
      </div>

      <div class="vtuber__content-item">
        <label for="lang-code">Region/Language Code</label>
        <input type="text" id="lang-code" name="lang-code" autocomplete="off" />
        <small class="description"
          >Minimum 2 characters (Example: <b>EN</b>)</small
        >
        <small class="error"></small>
      </div>

      <div class="vtuber__content-item">
        <label for="youtube-id">YouTube channel ID</label>
        <input
          type="text"
          id="youtube-id"
          name="youtube-id"
          autocomplete="off"
          ref="ytid"
        />
        <small class="description"
          >You can find the ID in the URL
          (https://www.youtube.com/channel/<b>HERE</b>) <br />(Example:
          <b>UCCzUftO8KOVkV4wQG1vkUvg</b>)</small
        >
        <small class="error"></small>
      </div>

      <div class="vtuber__content-item">
        <label for="nickname-twitch">Twitch Nickname</label>
        <input
          type="text"
          id="nickname-twitch"
          name="nickname-twitch"
          autocomplete="off"
          ref="twitchname"
        />
        <small class="description">
          You can find the Nickname in the URL (https://twitch.tv/<b>HERE</b>)
          <br />(Example: <b>ironmouse</b>)
        </small>
        <small class="error"></small>
      </div>

      <div class="vtuber__content-item">
        <label for="space-id">Bilibili Space ID</label>
        <input
          type="text"
          id="space-id"
          name="space-id"
          autocomplete="off"
          ref="biliid"
        />
        <small class="description">
          You can find the ID in the URL
          (https://space.bilibili.com/<b>HERE</b>) (not work in Bstation)
          <br />(Example: <b>339567211</b>)
        </small>
        <small class="error"></small>
      </div>

      <div class="vtuber__content-item">
        <label for="live-id">Bilibili Live ID</label>
        <input
          type="text"
          id="live-id"
          name="live-id"
          autocomplete="off"
          ref="liveid"
        />
        <small class="description">
          You can find the ID in the URL (https://live.bilibili.com/<b>HERE</b>)
          <br />(not work in Bstation) (Live ID and Space ID is not the same)
          <br />(Example: <b>14275133</b>)
        </small>
        <small class="error"></small>
      </div>

      <div class="vtuber__content-item">
        <label for="live-id">Bilibili Fanart</label>
        <input
          type="text"
          id="bili-art"
          name="bili-art"
          autocomplete="off"
          ref="liveid"
        />
        <small class="description">
         (Optional)
        </small>
        <small class="error"></small>
      </div>

      <div class="vtuber__content-item">
        <label for="twitter-username">Twitter Username</label>
        <input
          type="text"
          id="twitter-username"
          name="twitter-username"
          autocomplete="off"
        />
        <small class="description">
          (Optional) You can find the Username in the URL
          (https://twitter.com/<b>HERE</b>)
          <br />(Example: <b>Hana_Macchia</b>)
        </small>
        <small class="error"></small>
      </div>

      <div class="vtuber__content-item">
        <label for="fanart-hashtag">Twitter Fanart Hashtag</label>
        <input
          type="text"
          id="fanart-hashtag"
          name="fanart-hashtag"
          autocomplete="off"
        />
        <small class="description">
          (Optional) Get fanart from twitter. Find in social media vtuber or
          <a
            href="https://virtualyoutuber.fandom.com/"
            target="_blank"
            rel="noopener noreferrer"
            >wiki</a
          >. <br />(Example: <b>#ioarts</b> or <b>#しいなーと</b>)
        </small>
        <small class="error"></small>
      </div>

      <div class="vtuber__content-item">
        <label for="lewd-hashtag">Twitter NSFW fanart Hashtag</label>
        <input
          type="text"
          id="lewd-hashtag"
          name="lewd-hashtag"
          autocomplete="off"
        />
        <small class="description">
          (Optional) Get hnt4i fanart from twitter. Find in social media vtuber
          or
          <a
            href="https://virtualyoutuber.fandom.com/"
            target="_blank"
            rel="noopener noreferrer"
            >wiki</a
          >.</small
        >
        <small class="error"></small>
      </div>
    </div>
  </div>
</template>

<script>
import VtuberPlatform from "./VtuberPlatform_script.js"

export default { ...VtuberPlatform }
</script>

<style lang="scss" scoped>
.vtuber {
  &-form {
    @apply mb-2 rounded-lg bg-slate-200 md:ml-2;

    &.show {
      .arrow {
        @apply -rotate-90;
      }
      .vtuber__content {
        @apply h-[var(--contentHeight)] scale-y-100;
      }
    }

    &.errors {
      @apply bg-red-200;

      .vtuber-link {
        @apply bg-red-300;
      }
    }
  }
  &-link {
    @apply flex items-center justify-between rounded-lg bg-slate-300 px-4 py-2 shadow-md;

    &__icon {
      @apply flex items-center space-x-4;
    }
  }

  &__content {
    transition-property: "transform, height";
    @apply flex h-0 origin-top scale-y-0 flex-col duration-300 ease-in-out;

    &-item {
      @apply mx-2 my-1 flex flex-col first:mt-2;

      &.has-error {
        input {
          @apply bg-red-400 text-white;
        }

        .description {
          @apply hidden;
        }

        .error {
          @apply block;
        }
      }

      label {
        @apply ml-1;
      }

      input {
        @apply my-1 -translate-y-0.5 rounded-lg bg-slate-100 p-2 shadow-md transition duration-200 ease-in-out hover:translate-y-0 hover:shadow-sm focus:translate-y-0.5 focus:shadow-none focus:outline-none;
      }

      small {
        @apply text-xs;
      }

      .description {
        @apply text-gray-600;

        a {
          @apply text-blue-500 hover:text-blue-700;
        }
      }

      .error {
        @apply hidden text-red-500;
      }
    }
  }
}

.arrow {
  @apply transition-transform duration-300 ease-in-out;
}

.delete-vtuber {
  @apply rounded-full bg-white p-1 px-1.5;

  &.one {
    @apply bg-slate-200 opacity-70;
  }

  &.confirm {
    @apply bg-red-500 text-white;
  }
}
</style>
