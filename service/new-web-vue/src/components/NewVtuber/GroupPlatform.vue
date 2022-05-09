<template>
  <div class="platform-group show" ref="platform">
    <a
      class="platform-link"
      href="#"
      onclick="return false"
      @click="toggleContent"
    >
      <span class="platform-link__text">{{
        platform == "youtube" ? "YouTube" : "BiliBili"
      }}</span>
      <span class="platform-link__icon">
        <a
          class="delete-platform"
          href="#"
          onclick="return false"
          @click="deletePlatform"
          ><font-awesome-icon icon="trash-can" class="fa-fw"
        /></a>
        <font-awesome-icon icon="chevron-down" class="arrow" />
      </span>
    </a>
    <div class="platform-group__content" ref="content">
      <div class="platform-group__content-item">
        <label for="type-platform">Type Platform</label>
        <select
          id="type-platform"
          name="type-platform"
          @change="platform = $event.target.value"
          :value="platform"
        >
          <option value="youtube">YouTube</option>
          <option value="bilibili">BiliBili</option>
        </select>
      </div>
      <div class="platform-group__content-item" v-if="platform === `youtube`">
        <label for="youtube-id">YouTube ID</label>
        <input type="text" id="youtube-id" name="youtube-id" />
        <small class="description"
          >You can find the ID in the URL <br />(Example:
          https://www.youtube.com/channel/<b>UCCzUftO8KOVkV4wQG1vkUvg</b>)</small
        >
      </div>
      <div class="platform-group__content-item" v-if="platform === `bilibili`">
        <label for="space-id">Bilibili Space ID</label>
        <input type="text" id="space-id" name="space-id" />
        <small class="description">
          You can find the ID in the URL (not work in Bstation)
          <br />(Example: https://space.bilibili.com/<b>339567211</b>)
        </small>
      </div>
      <div class="platform-group__content-item" v-if="platform === `bilibili`">
        <label for="live-id">Bilibili Live ID</label>
        <input type="text" id="live-id" name="live-id" />
        <small class="description">
          You can find the ID in the URL (not work in Bstation)
          <br />(Example: https://live.bilibili.com/<b>14275133</b>)
        </small>
      </div>
      <div class="platform-group__content-item">
        <label for="lang-code">Region/Language Code</label>
        <input type="text" id="lang-code" name="lang-code" />
        <small class="description"
          >Minimum 2 characters (Example: <b>EN</b>)</small
        >
      </div>
    </div>
  </div>
</template>
<script>
import { library } from "@fortawesome/fontawesome-svg-core"
import { faChevronDown, faTrashCan } from "@fortawesome/free-solid-svg-icons"

library.add(faChevronDown, faTrashCan)

export default {
  data() {
    return {
      platform: "youtube",
    }
  },
  props: {
    id: {
      type: Number,
      required: true,
    },
    set: {
      type: String,
    },
  },
  mounted() {
    const platformGroups = document.querySelectorAll(".platform-group")

    const platformContent = this.$refs.content
    const platformGroup = this.$refs.platform
    // set style variable (--contentHeight)
    platformContent.style.setProperty(
      "--contentHeight",
      `${this.platform === "youtube" ? "286" : "398"}px`
    )

    const filteredGroups = [...platformGroups].filter(
      (c) => c !== platformGroup
    )

    filteredGroups.forEach((c) => c.classList.remove("show"))

    this.$watch(
      () => this.platform,
      () => {
        platformContent.style.setProperty(
          "--contentHeight",
          `${this.platform === "youtube" ? "286" : "398"}px`
        )
      }
    )

    document.body.addEventListener("click", (e) => {
      if (!e.target.closest(".delete-platform")) {
        document.querySelectorAll(".delete-platform").forEach((platform) => {
          platform.classList.remove("confirm")
        })
      }
    })

    this.$watch(
      () => this.set,
      () => (this.platform = this.set)
    )
  },
  methods: {
    deletePlatform(e) {
      document.querySelectorAll(".delete-platform").forEach((platform) => {
        if (platform !== e.target.closest(".delete-platform"))
          platform.classList.remove("confirm")
      })
      if (e.target.closest(".delete-platform").classList.contains("confirm"))
        this.$emit("delete", this.id)
      else e.target.closest(".delete-platform").classList.add("confirm")
    },
    toggleContent(e) {
      if (e.target.closest(".delete-platform")) return
      const group = e.target.closest(".platform-group")

      const groups = [...document.querySelectorAll(".platform-group")].filter(
        (c) => c !== group
      )

      groups.forEach((c) => c.classList.remove("show"))

      if (group.classList.contains("show")) group.classList.remove("show")
      else group.classList.add("show")
    },
  },
}
</script>

<style lang="scss" scoped>
.platform {
  &-group {
    @apply mb-2 rounded-lg bg-slate-200 md:ml-2;

    &.show {
      .arrow {
        @apply -rotate-90;
      }
      .platform-group__content {
        @apply h-[var(--contentHeight)] scale-y-100;
      }
    }

    &__content {
      transition-property: transform, height;
      @apply flex  h-0 origin-top scale-y-0 flex-col duration-300 ease-in-out;

      &-item {
        @apply mx-2 my-1 flex flex-col;

        label {
          @apply ml-1;
        }

        select {
          @apply my-1 -translate-y-0.5 rounded-lg bg-slate-100 p-2 shadow-md transition duration-200 ease-in-out hover:translate-y-0 hover:shadow-sm focus:translate-y-0.5 focus:shadow-none;
        }

        input {
          @apply my-1 -translate-y-0.5 rounded-lg bg-slate-100 p-2 shadow-md transition duration-200 ease-in-out hover:translate-y-0 hover:shadow-sm focus:translate-y-0.5 focus:shadow-none focus:outline-none;
        }

        small {
          @apply text-xs;
        }

        .description {
          @apply text-gray-600;
        }

        .error {
          @apply text-red-500;
        }
      }
    }
  }

  &-link {
    @apply flex items-center justify-between rounded-lg bg-slate-300 px-4 py-2 shadow-md;

    &__icon {
      @apply flex items-center space-x-4;
    }
  }
}

.arrow {
  @apply transition-transform duration-300 ease-in-out;
}

.delete-platform {
  @apply rounded-full bg-white p-1 px-1.5;

  &.confirm {
    @apply bg-red-500 text-white;
  }
}
</style>
