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
        <label for="youtube-id">YouTube channel ID</label>
        <input
          type="text"
          id="youtube-id"
          name="youtube-id"
          autocomplete="off"
        />
        <small class="description"
          >You can find the ID in the URL
          (https://www.youtube.com/channel/<b>HERE</b>) <br />(Example:
          <b>UCCzUftO8KOVkV4wQG1vkUvg</b>)</small
        >
        <small class="error"></small>
      </div>

      <div class="platform-group__content-item" v-if="platform === `bilibili`">
        <label for="space-id">Bilibili Space ID</label>
        <input type="text" id="space-id" name="space-id" autocomplete="off" />
        <small class="description">
          You can find the ID in the URL
          (https://space.bilibili.com/<b>HERE</b>) (not work in Bstation)
          <br />(Example: <b>339567211</b>)
        </small>
        <small class="error"></small>
      </div>

      <div class="platform-group__content-item" v-if="platform === `bilibili`">
        <label for="live-id">Bilibili Live ID</label>
        <input type="text" id="live-id" name="live-id" autocomplete="off" />
        <small class="description">
          You can find the ID in the URL (https://live.bilibili.com/<b>HERE</b>)
          <br />(not work in Bstation) (Live ID and Space ID is not the same)
          <br />(Example: <b>14275133</b>)
        </small>
        <small class="error"></small>
      </div>

      <div class="platform-group__content-item">
        <label for="lang-code">Region/Language Code</label>
        <input type="text" id="lang-code" name="lang-code" autocomplete="off" />
        <small class="description"
          >Minimum 2 characters (Example: <b>EN</b>)</small
        >
        <small class="error"></small>
      </div>
    </div>
  </div>
</template>
<script>
import { library } from "@fortawesome/fontawesome-svg-core"
import { faChevronDown, faTrashCan } from "@fortawesome/free-solid-svg-icons"

library.add(faChevronDown, faTrashCan)

import trim from "validator/lib/trim"

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
  emits: ["delete", "error"],
  async mounted() {
    this.$refs.platform
      .querySelectorAll(".platform-group__content-item")
      .forEach((item) => {
        if (item.querySelector("select")) return
        item.querySelector("input").value = ""
        item.classList.remove("has-error")
        item.querySelector(".error").innerHTML = ""
      })

    await this.checkHeight()

    const platform = this.$refs.platform

    platform.addEventListener("input", async (e) => {
      e.target.parentElement.classList.toggle(
        "has-error",
        !(await this.checkFilled(e.target))
      )

      await new Promise((resolve) => setTimeout(resolve, 60))
      await this.checkAllFilled(e.target)
      this.checkHeight()
    })

    let activeElement = null

    document.body.addEventListener("click", async (e) => {
      if (
        activeElement &&
        e.target !== activeElement &&
        activeElement.tagName === "INPUT" &&
        activeElement.closest(".platform-group__content")
      ) {
        activeElement.parentElement.classList.toggle(
          "has-error",
          !(await this.checkFilled(activeElement))
        )

        await new Promise((resolve) => setTimeout(resolve, 60))
        await this.checkAllFilled(activeElement)
        this.checkHeight()
      }

      activeElement = e.target
    })

    this.$watch(
      () => this.id,
      async () => {
        const inputs = platform.querySelectorAll("input")

        inputs.forEach(async (input) =>
          input.parentElement.classList.toggle(
            "has-error",
            !(await this.checkFilled(activeElement))
          )
        )

        await new Promise((resolve) => setTimeout(resolve, 60))
        await this.checkAllFilled(null)
        this.checkHeight()
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
    async checkFilled(element) {
      const notInput = element.tagName !== "INPUT"
      // this.calculateHeight(element)
      if (notInput) return true
      return await this.checkValidate(element)
    },
    async checkAllFilled(e) {
      const platform = this.$refs.platform
      if (!platform) return

      const inputs = platform.querySelectorAll("input")
      let count = 0
      for (const input of inputs) {
        const validate = await this.checkValidate(input)
        if (!validate) break
        count++
      }

      const error = count < inputs.length && e.tagName === "INPUT"
      platform.classList.toggle("errors", error)
      this.$emit("error", { id: this.id, error })
    },
    async checkValidate(e) {
      const youtubeId = e.name === "youtube-id"
      const biliSpaceId = e.name === "space-id"
      const biliLiveId = e.name === "live-id"
      const langCode = e.name === "lang-code"

      const errorText = e.parentElement.querySelector(".error")
      const value = trim(e.value)

      if ((youtubeId || biliSpaceId || biliLiveId || langCode) && !value) {
        errorText.innerText = "This field is required"
        return false
      }

      // check is valid youtube id

      const isChannelId = value.match(/^UC[a-zA-Z0-9-_]{22}$/)

      if (youtubeId && !isChannelId) {
        errorText.innerText = "This is not a valid channel ID, find inside URL"
        return false
      }

      // check is valid bilibili live id or space id

      const isBiliBiliId = value.match(/^\d+$/)

      if ((biliLiveId || biliSpaceId) && !isBiliBiliId) {
        errorText.innerText = "This is not a valid live ID, find inside URL"
        return false
      }

      // check region/language code
      const isRegionCode = value.match(/^[a-zA-Z]{2}$/)

      if (langCode && !isRegionCode) {
        errorText.innerText = "Please enter a valid region/language code"
        return false
      }

      return true
    },
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

      group.classList.toggle("show")
    },
    async checkHeight() {
      if (!this.$refs.content) return
      await new Promise((resolve) => setTimeout(resolve, 60))
      const calculatedHeight = [...this.$refs.content.children].reduce(
        (t, c) => {
          // get margin and padding
          const margin =
            parseInt(getComputedStyle(c).marginTop) +
            parseInt(getComputedStyle(c).marginBottom)
          const padding =
            parseInt(getComputedStyle(c).paddingTop) +
            parseInt(getComputedStyle(c).paddingBottom)

          // get total height
          return t + c.offsetHeight + margin + padding
        },
        0
      )

      this.$refs.content.style.setProperty(
        "--contentHeight",
        `${calculatedHeight}px`
      )
    },
  },
}
</script>

<style lang="scss" scoped>
.platform {
  &-group {
    @apply mb-2 rounded-lg bg-slate-200 dark:bg-slate-500 md:ml-2;

    &.show {
      .arrow {
        @apply -rotate-90;
      }
      .platform-group__content {
        @apply h-[var(--contentHeight)] scale-y-100;
      }
    }

    &.errors .platform-link {
        @apply bg-red-300 dark:bg-red-600;
      }
    }

    &__content {
      transition-property: "transform, height";
      @apply flex h-0 origin-top scale-y-0 flex-col duration-300 ease-in-out;

      &-item {
        @apply mx-2 my-1 flex flex-col first:mt-2;

        &.has-error {
          input {
            @apply bg-red-400 dark:bg-red-600;
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

        select {
          @apply my-1 -translate-y-0.5 rounded-lg bg-slate-100 p-2 shadow-md transition duration-200 ease-in-out hover:translate-y-0 hover:shadow-sm focus:translate-y-0.5 focus:shadow-none dark:bg-slate-600;
        }

        input {
          @apply my-1 -translate-y-0.5 rounded-lg bg-slate-100 p-2 shadow-md transition duration-200 ease-in-out hover:translate-y-0 hover:shadow-sm focus:translate-y-0.5 focus:shadow-none focus:outline-none dark:bg-slate-600;
        }

        small {
          @apply text-xs;
        }

        .description {
          @apply text-gray-600 dark:text-slate-50;
        }

        .error {
          @apply hidden text-red-500;
        }
      }
    }
  }

  &-link {
    @apply flex items-center justify-between rounded-lg bg-slate-300 px-4 py-2 shadow-md dark:bg-slate-600;

    &__icon {
      @apply flex items-center space-x-4;
    }
  }
}

.arrow {
  @apply transition-transform duration-300 ease-in-out;
}

.delete-platform {
  @apply rounded-full bg-white p-1 px-1.5 transition duration-200 ease-in-out dark:text-slate-700;

  &.confirm {
    @apply bg-red-500 text-white;
  }
}
</style>
