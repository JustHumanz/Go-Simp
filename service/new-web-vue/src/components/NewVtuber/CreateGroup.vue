<script setup>
import GroupPlatform from "./GroupPlatform.vue"
</script>

<template>
  <h1 class="title-req-group">Request new Group</h1>
  <form @submit="sendGroup" class="form-req-group">
    <div class="form-group">
      <label for="group-name">Group Name</label>
      <input
        type="text"
        id="group-name"
        name="group-name"
        class="form-control"
        autocomplete="off"
        placeholder="Group Name"
      />
      <small class="error"></small>
    </div>
    <div class="form-group">
      <label for="group-icon">URL Icon</label>
      <input
        type="text"
        id="group-icon"
        name="group-icon"
        class="form-control"
        autocomplete="off"
        placeholder="URL Icon"
      />
      <small class="description"
        >(Optional) It is recommended to fill this in accordingly to make it
        easier to find the icon.</small
      >
      <small class="error"></small>
    </div>
    <div class="form-group">
      <button type="button" name="add-platform" @click="addPlatform">
        Add Platform
      </button>
    </div>
    <div class="platforms">
      <group-platform
        v-for="platform in platforms"
        :id="platform.id"
        :set="platform.set"
        @delete="deletePlatform"
      />
    </div>
    <!-- Submit and back btn -->
    <div class="form-group-btn">
      <button type="submit" class="submit">Submit</button>
      <button type="button" class="back" @click="$emit(`back`, true)">
        Back
      </button>
    </div>
  </form>
</template>

<script>
import trim from "validator/lib/trim"
import isUrl from "validator/lib/isURL"
export default {
  data() {
    return {
      platforms: [],
    }
  },
  props: {
    groups: {
      type: Number,
      default: [],
    },
  },
  mounted() {
    console.log(this.groups)
    const submitBtn = document.querySelector(".submit")

    submitBtn.disabled = true

    document.body.addEventListener("input", (e) => {
      e.target.parentElement.classList.toggle(
        "has-error",
        !this.checkFilled(e.target)
      )

      this.checkAllFilled()
    })

    let activeElement = null

    document.body.addEventListener("click", (e) => {
      if (activeElement && e.target !== activeElement) {
        activeElement.parentElement.classList.toggle(
          "has-error",
          !this.checkFilled(activeElement)
        )
      }

      activeElement = e.target
    })
  },
  methods: {
    sendGroup(e) {
      e.preventDefault()
      console.log(e.target)
    },
    async addPlatform(e) {
      if (this.platforms.length == 9) e.target.disabled = true
      this.platforms.push({
        id: this.platforms.length,
        set: "youtube",
      })

      await new Promise((resolve) => setTimeout(resolve, 60))
      this.checkAllFilled()
    },
    async deletePlatform(id) {
      const addPlatformBtn = document.querySelector("[name=add-platform]")
      if (addPlatformBtn.disabled) addPlatformBtn.disabled = false

      const platforms = [...document.querySelector(".platforms").children]
      // remove confirm class from delete-platform
      platforms.forEach((platform) => {
        platform.querySelector(".delete-platform").classList.remove("confirm")
      })

      const filteredPlatforms = platforms.filter((platform, index) => {
        return index !== id
      })
      this.platforms.splice(id, 1)

      await new Promise((resolve) => setTimeout(resolve, 60))
      this.checkAllFilled()

      if (!filteredPlatforms.length) return

      this.platforms.map((platform, index) => {
        platform.id = index
      })

      filteredPlatforms.forEach((platform, index) => {
        const select = platform.querySelector("select")
        this.platforms[index].set = select.value
      })

      const newPlatforms = [...document.querySelector(".platforms").children]

      newPlatforms.forEach((platform, index) => {
        const input = platform.querySelectorAll("input")
        const oldInput = filteredPlatforms[index].querySelectorAll("input")

        platform.classList = filteredPlatforms[index].classList

        input.forEach((inp, i) => (inp.value = oldInput[i].value))
      })
    },
    checkFilled(element) {
      const notInput = element.tagName !== "INPUT"
      const iconUrl = element.name === "group-icon"
      const errorText = element.parentElement.querySelector(".error")

      this.calculateHeight(element)

      if (notInput) return true
      const value = trim(element.value)

      // toggle error class in parent element
      if (!iconUrl && !value) {
        errorText.innerText = "This field is required"
        return false
      }

      // when group name same
      if (element.name === "group-name") {
        const groupNameExist = this.groups.some(
          (group) =>
            group.GroupName.toLowerCase().replace("_", " ") ===
            value.toLowerCase()
        )

        if (groupNameExist) {
          errorText.innerText = "Group name already exist, please find new one"
          return false
        }
      }

      // when is url
      if (iconUrl) {
        if (value && !isUrl(value)) {
          errorText.innerText = "Please enter a valid URL (image)"
          return false
        }
      }

      // check is valid youtube id
      if (element.name === "youtube-id") {
        const isChannelId = value.match(/^UC[a-zA-Z0-9-_]{22}$/)

        if (!isChannelId) {
          errorText.innerText =
            "This is not a valid channel ID, find inside URL"
          return false
        }
      }

      // check is valid bilibili live id or space id
      if (element.name === "live-id") {
        const isBiliBiliLiveId = value.match(/^\d+$/)

        if (!isBiliBiliLiveId) {
          errorText.innerText = "This is not a valid live ID, find inside URL"
          return false
        }
      }

      if (element.name === "space-id") {
        const isBiliBiliSpaceId = value.match(/^\d+$/)

        if (!isBiliBiliSpaceId) {
          errorText.innerText = "This is not a valid space ID, find inside URL"
          return false
        }
      }

      // check region/language code
      if (element.name === "lang-code") {
        const isRegionCode = value.match(/^[a-zA-Z]{2}$/)

        if (!isRegionCode) {
          errorText.innerText = "Please enter a valid region/language code"
          return false
        }
      }
      return true
    },
    checkAllFilled() {
      const input = document.querySelectorAll("input")
      console.log(input)
    },
    async calculateHeight(element) {
      const outsideDiv = element.parentElement.parentElement

      if (!outsideDiv.classList.contains("platform-group__content")) return
      await new Promise((resolve) => setTimeout(resolve, 60))

      const calculatedHeight = [...outsideDiv.children].reduce((t, c) => {
        // get margin and padding
        const margin =
          parseInt(getComputedStyle(c).marginTop) +
          parseInt(getComputedStyle(c).marginBottom)
        const padding =
          parseInt(getComputedStyle(c).paddingTop) +
          parseInt(getComputedStyle(c).paddingBottom)

        // get total height
        return t + c.offsetHeight + margin + padding
      }, 0)

      outsideDiv.style.setProperty("--contentHeight", `${calculatedHeight}px`)
    },
  },
}
</script>
<style lang="scss" scoped>
.title-req-group {
  @apply my-2 text-lg font-semibold;
}

.form-req-group {
  @apply flex w-full flex-col space-y-2 md:w-3/6;
}

.form-group {
  @apply flex flex-col md:ml-2;

  &.has-error {
    input {
      @apply bg-red-200;
    }

    .description {
      @apply hidden;
    }

    .error {
      @apply block;
    }
  }

  input {
    @apply my-1 -translate-y-0.5 rounded-lg bg-slate-200 p-2 shadow-md transition duration-200 ease-in-out hover:translate-y-0 hover:shadow-sm focus:translate-y-0.5 focus:shadow-none focus:outline-none;
  }

  small {
    @apply text-xs;
  }

  .description {
    @apply text-gray-600;
  }

  .error {
    @apply hidden text-red-500;
  }

  button {
    @apply -translate-y-0.5 rounded-lg bg-blue-400 p-2 font-semibold text-white shadow-md transition duration-200 ease-in-out hover:translate-y-0 hover:shadow-sm disabled:translate-y-0.5 disabled:opacity-70 disabled:shadow-none;
  }
}

.form-group-btn {
  @apply grid grid-rows-2 gap-2 pb-4 md:grid-cols-2 md:grid-rows-none;

  button {
    @apply -translate-y-0.5 rounded-lg p-2 font-semibold text-white shadow-md transition duration-200 ease-in-out first:bg-green-500 last:bg-red-600 hover:translate-y-0 hover:shadow-sm disabled:translate-y-0 disabled:opacity-70 disabled:shadow-none;
  }

  .platforms {
    @apply flex flex-col gap-2;
  }
}
</style>

<style lang="scss">
.platforms {
  .platform {
    @apply flex flex-col md:ml-2;

    .delete-btn {
      @apply mt-2 -translate-y-0.5 rounded-lg bg-red-400 px-2 py-1 text-sm font-semibold text-white shadow-md transition duration-200 ease-in-out hover:translate-y-0 hover:shadow-sm disabled:translate-y-0.5 disabled:opacity-70 disabled:shadow-none;
    }

    label {
      @apply text-lg;
    }

    select {
      @apply my-1 -translate-y-0.5 rounded-lg bg-slate-200 p-2 shadow-md transition duration-200 ease-in-out hover:translate-y-0 hover:shadow-sm focus:translate-y-0.5 focus:shadow-none focus:outline-none;
    }

    input {
      @apply my-1 -translate-y-0.5 rounded-lg bg-slate-200 p-2 shadow-md transition duration-200 ease-in-out hover:translate-y-0 hover:shadow-sm focus:translate-y-0.5 focus:shadow-none focus:outline-none;
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
</style>
