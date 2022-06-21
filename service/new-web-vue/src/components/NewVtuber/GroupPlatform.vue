<template>
  <div class="platform-group show" ref="platform" :data-id="id">
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
        <!-- <select
          id="type-platform"
          name="type-platform"
          @change="platform = $event.target.value"
          :value="platform"
        >
          <option value="youtube">YouTube</option>
          <option value="bilibili">BiliBili</option>
        </select> -->
        <div class="custom-select" style="--maxHeight: 80px; --minHeight: 40px">
          <div class="custom-select-area">
            <a
              class="selected-item"
              href="#"
              onclick="return false"
              v-if="platform == 'youtube'"
            >
              <font-awesome-icon :icon="['fab', 'youtube']" class="fa-fw" />
              <span class="selected-item__span">YouTube</span>
            </a>
            <a
              class="selected-item"
              href="#"
              onclick="return false"
              v-if="platform == 'bilibili'"
            >
              <font-awesome-icon :icon="['fab', 'bilibili']" class="fa-fw" />
              <span class="selected-item__span">BiliBili</span>
            </a>
            <div class="select-items">
              <!-- YouTube -->
              <input
                type="radio"
                class="item__radio"
                name="type-platform"
                id="youtube"
                value="youtube"
                checked
              />
              <label for="youtube" class="item__label">
                <font-awesome-icon :icon="['fab', 'youtube']" class="fa-fw" />
                <span class="item__label-span">YouTube</span>
              </label>
              <!-- BiliBili -->
              <input
                type="radio"
                class="item__radio"
                name="type-platform"
                id="bilibili"
                value="bilibili"
              />
              <label for="bilibili" class="item__label">
                <font-awesome-icon :icon="['fab', 'bilibili']" class="fa-fw" />
                <span class="item__label-span">BiliBili</span>
              </label>
            </div>
          </div>
        </div>
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
        <label for="lang-code">Region/Language</label>

        <div class="region-area" v-if="toggleLang === id">
          <div class="regions">
            <div class="region" v-for="region in regions">
              <input
                type="radio"
                name="region"
                class="region__radio"
                :id="region.code"
                :value="region.code"
              />
              <label
                :for="region.code"
                class="region__label"
                @mousedown="setReg"
              >
                <img
                  :src="`/assets/flags/${region.code}.svg`"
                  onerror="this.src='/assets/flags/none.svg'"
                  class="region__label-flag"
                />
                <span class="region__label-span">{{ region.name }}</span>
              </label>
            </div>
            <div class="region" v-if="!regions.length">
              <label class="region__label disabled">
                <span>No Regions Found</span>
              </label>
            </div>
          </div>
        </div>

        <input
          type="text"
          id="lang-code"
          name="lang-code"
          autocomplete="off"
          @click="openLang"
          @blur="toggleLang = null"
          @input="findReg"
        />
        <small class="description">Find Available Regions/Languages</small>
        <small class="error"></small>
      </div>
    </div>
  </div>
</template>
<script>
import { library } from "@fortawesome/fontawesome-svg-core"
import { faChevronDown, faTrashCan } from "@fortawesome/free-solid-svg-icons"
import { faYoutube, faBilibili } from "@fortawesome/free-brands-svg-icons"

library.add(faChevronDown, faTrashCan, faYoutube, faBilibili)

import Regions from "@/regions.json"

import trim from "validator/lib/trim"

export default {
  data() {
    return {
      platform: "youtube",
      toggleLang: null,
      searchLang: "",
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
  emits: ["delete", "error", "set"],
  async mounted() {
    const platformForm = this.$refs.platform
    platformForm
      .querySelectorAll(".platform-group__content-item")
      .forEach((item) => {
        if (item.querySelector(".custom-select")) return
        item.querySelector("input").value = ""
        item.classList.remove("has-error")
        item.querySelector(".error").innerHTML = ""
      })

    await this.checkHeight()

    platformForm.addEventListener("input", async (e) => {
      e.target.parentElement.classList.toggle(
        "has-error",
        !(await this.checkFilled(e.target, true))
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
        activeElement?.closest(".platform-group") === platformForm
      ) {
        activeElement.parentElement.classList.toggle(
          "has-error",
          !(await this.checkFilled(activeElement))
        )

        await new Promise((resolve) => setTimeout(resolve, 60))
        await this.checkAllFilled(activeElement)
      }

      activeElement = e.target
      this.checkHeight()
    })

    platformForm.addEventListener("mousedown", async (e) => {
      const labelSelected = e.target.closest(".item__label")

      if (labelSelected) {
        const platforms = [...platformForm.querySelectorAll(".item__radio")]
        // check label connected to radio
        platforms.forEach((item) => {
          const selected = item.nextElementSibling === labelSelected
          if (!selected) return

          item.checked = true
          this.platform = item.value
          this.$emit("set", { id: this.id, set: item.value })
        })
      }

      this.checkHeight()
    })

    document.body.addEventListener("click", (e) => {
      if (!e.target.closest(".delete-platform")) {
        document.querySelectorAll(".delete-platform").forEach((platform) => {
          platform.classList.remove("confirm")
        })
      }
    })

    this.$watch(
      () => this.set,
      () => (this.platform = this.set),
      { immediate: true }
    )
  },

  computed: {
    regions() {
      return Regions.filter(
        (region) =>
          region.code.toLowerCase().includes(this.searchLang.toLowerCase()) ||
          region.name.toLowerCase().includes(this.searchLang.toLowerCase())
      )
    },
  },
  methods: {
    async checkFilled(element, isInput = false) {
      const notInput = element.tagName !== "INPUT"
      // this.calculateHeight(element)
      if (notInput || isInput) return true
      return await this.checkValidate(element, isInput)
    },
    async checkAllFilled(e) {
      const platform = this.$refs.platform
      if (!platform) return

      const inputs = platform.querySelectorAll("input")
      let count = 0
      for (const input of inputs) {
        const validate = await this.checkValidate(input, true)
        if (!validate) break
        count++
      }

      const error = count < inputs.length && e.tagName === "INPUT"
      platform.classList.toggle("errors", error)
      this.$emit("error", { id: this.id, error })
    },
    async checkValidate(e, isInput = false) {
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
      const isRegion = Regions.find((region) => region.name === value)

      if (langCode && !isRegion) {
        if (!isInput) errorText.innerText = "Please enter a valid region"
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
      await new Promise((resolve) => setTimeout(resolve, 60))
      if (!this.$refs.content) return
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

    openLang(e) {
      const input = e.target

      if (input.value && this.searchLang !== "")
        input.setSelectionRange(0, input.value.length)

      this.toggleLang = this.id
      this.searchLang = ""
    },

    setReg(e) {
      // get value radio
      const selectedRegion =
        e.target.closest(".region__label")?.previousElementSibling.value

      // set value to input
      const legionInput = e.target
        .closest(".platform-group__content-item")
        .querySelector("input[name='lang-code']")
      legionInput.value = Regions.find((r) => r.code === selectedRegion)?.name
      this.searchLang = Regions.find((r) => r.code === selectedRegion)?.name
      this.checkAllFilled(e.target)
    },
    findReg(e) {
      e.target.classList.toggle(
        "region-selected",
        !!this.regions.find((r) => r.name === e.target.value)
      )
      this.searchLang = e.target.value
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

    &__content {
      transition-property: "transform, height";
      @apply flex h-0 origin-top scale-y-0 flex-col overflow-hidden duration-300 ease-in-out;

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

        .custom-select {
          @apply relative h-[40px];

          &-area {
            transition-property: height transform shadow;
            @apply z-[4] my-1 h-[40px] w-full -translate-y-0.5 cursor-pointer overflow-hidden rounded-lg bg-slate-100 shadow-md duration-200 ease-in-out focus-within:absolute focus-within:h-[80px] focus-within:translate-y-0.5 hover:translate-y-0 hover:shadow-sm focus-within:hover:translate-y-0.5  focus-within:hover:shadow-md dark:bg-slate-600;
            .selected-item {
              @apply inline-block w-full p-2 focus:hidden;

              &__span {
                @apply ml-2;
              }
            }

            .select-items {
              @apply flex flex-col;

              input {
                @apply hidden;
              }

              .item__label {
                @apply m-0 cursor-pointer p-2 hover:bg-slate-300 dark:hover:bg-slate-700;

                &-span {
                  @apply ml-2;
                }
              }
            }
          }
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

.regions {
  @apply absolute bottom-0 mb-1 flex max-h-48 w-full flex-col overflow-hidden overflow-y-scroll rounded-md bg-slate-100 shadow-md dark:bg-slate-600;

  .region {
    @apply flex;

    &__radio {
      @apply hidden;
    }

    &__label {
      @apply m-0 flex h-full w-full cursor-pointer p-2 hover:bg-slate-300 dark:hover:bg-slate-700;

      &.disabled {
        @apply cursor-not-allowed hover:bg-transparent dark:hover:bg-transparent;
      }

      &-flag {
        @apply h-6 w-6 rounded-md object-contain drop-shadow-md;
      }

      &-span {
        @apply ml-2 text-center;
      }
    }
  }
}

.region-area {
  @apply relative w-full;
}
</style>
