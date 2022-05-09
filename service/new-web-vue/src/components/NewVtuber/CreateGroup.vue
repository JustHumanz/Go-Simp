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
        autocorrect="off"
        placeholder="Group Name"
      />
    </div>
    <div class="form-group">
      <label for="group-icon">URL Icon</label>
      <input
        type="text"
        id="group-icon"
        name="group-icon"
        class="form-control"
        autocorrect="off"
        placeholder="URL Icon"
      />
      <small class="description"
        >(Optional) It is recommended to fill this in accordingly to make it
        easier to find the icon.</small
      >
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
export default {
  data() {
    return {
      platforms: [],
    }
  },
  mounted() {
    const submitBtn = document.querySelector(".submit")

    submitBtn.disabled = true

    document.body.addEventListener("input", (e) => {
      this.checkFilled(e.target)
    })

    let activeElement = null

    document.body.addEventListener("click", (e) => {
      if (e.target !== activeElement) {
        this.checkFilled(e.target)
      }

      activeElement = e.target
    })
  },
  methods: {
    sendGroup(e) {
      e.preventDefault()
      console.log(e.target)
    },
    addPlatform(e) {
      if (this.platforms.length == 9) e.target.disabled = true
      this.platforms.push({
        id: this.platforms.length,
        set: "youtube",
      })
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

      if (!filteredPlatforms.length) return

      this.platforms.map((platform, index) => {
        platform.id = index
      })

      filteredPlatforms.forEach((platform, index) => {
        const select = platform.querySelector("select")
        this.platforms[index].set = select.value
      })

      await new Promise((resolve) => setTimeout(resolve, 60))
      const newPlatforms = [...document.querySelector(".platforms").children]

      newPlatforms.forEach((platform, index) => {
        const input = platform.querySelectorAll("input")
        const oldInput = filteredPlatforms[index].querySelectorAll("input")

        platform.classList = filteredPlatforms[index].classList

        input.forEach((inp, i) => (inp.value = oldInput[i].value))
      })
    },
    checkFilled(element) {},
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
