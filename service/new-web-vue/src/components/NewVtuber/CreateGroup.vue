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
    <div class="platforms"></div>
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
import {} from "validator"
export default {
  mounted() {
    document.body.addEventListener("click", (e) => {
      if (e.target.name === "del-platform") {
        // get div platform
        const divPlatform = e.target.closest(".platform")
        // get div platforms
        const platforms = divPlatform.closest(".platforms")
        // delete div platform
        divPlatform.remove()

        // check add platform disabled
        document.querySelector("[name=add-platform]").disabled = false

        // looping platforms
        platforms.childNodes.forEach((pl, index) => {
          // get label
          const label = pl.querySelector("label")
          // change for inside label
          label.setAttribute("for", `platform-${index + 1}`)

          // get select
          const select = pl.querySelector("select")
          // change name and id inside select
          select.setAttribute("name", `platform-${index + 1}`)
          select.setAttribute("id", `platform-${index + 1}`)
        })
      }
    })

    document.body.addEventListener("change", this.selectPlatform)
  },
  methods: {
    sendGroup(e) {
      e.preventDefault()
      console.log(e.target)
    },
    addPlatform(e) {
      const platforms = document.querySelector(".platforms")
      const scopeId = this.$options.__scopeId
      // check total childs
      const totalChilds = platforms.children.length
      if (totalChilds === 9) e.target.disabled = true
      // add div platform inside platforms
      const divPlatform = platforms.appendChild(document.createElement("div"))
      divPlatform.classList.add("platform")
      divPlatform.setAttribute(scopeId, "")

      // add label inside div platform
      const label = divPlatform.appendChild(document.createElement("label"))
      label.setAttribute("for", `platform-${totalChilds + 1}`)
      label.setAttribute(scopeId, "")
      label.innerText = "Platform"

      // add select with option youtube or bilibili inside div platform
      const select = divPlatform.appendChild(document.createElement("select"))
      select.setAttribute("name", `platform-${totalChilds + 1}`)
      select.setAttribute("id", `platform-${totalChilds + 1}`)
      select.setAttribute(scopeId, "")
      select.value = ""

      // add value null disabled inside select
      const option = select.appendChild(document.createElement("option"))
      option.setAttribute("value", "")
      option.setAttribute("disabled", "")
      option.innerText = "Select Platform"

      const optionYoutube = select.appendChild(document.createElement("option"))
      optionYoutube.setAttribute(scopeId, "")
      optionYoutube.setAttribute("value", "youtube")
      optionYoutube.innerText = "Youtube"

      const optionBilibili = select.appendChild(
        document.createElement("option")
      )
      optionBilibili.setAttribute(scopeId, "")
      optionBilibili.setAttribute("value", "bilibili")
      optionBilibili.innerText = "BiliBili"

      // add delete btn inside div platform
      const deleteBtn = divPlatform.appendChild(
        document.createElement("button")
      )
      deleteBtn.setAttribute(scopeId, "")
      deleteBtn.classList.add("delete-btn")
      deleteBtn.setAttribute("name", "del-platform")
      deleteBtn.setAttribute("type", "button")
      deleteBtn.innerText = "Delete"
    },
    selectPlatform(e) {
      // get element select
      if (
        !e.target.tagName === "SELECT" &&
        !e.target.name.includes("platform-")
      )
        return

      const divPlatform = e.target.closest(".platform")
      // delete input-platform-group when exist
      const existPlatformGroup = divPlatform.querySelector(
        ".input-platform-group"
      )
      if (existPlatformGroup) existPlatformGroup.remove()
      // create div input-platform-group
      const divInputPlatformGroup = divPlatform.appendChild(
        document.createElement("div")
      )
      divInputPlatformGroup.classList.add("input-platform-group")
      divInputPlatformGroup.innerText = e.target.value
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

  label {
    @apply text-lg;
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

  button {
    @apply -translate-y-0.5 rounded-lg bg-blue-400 p-2 font-semibold text-white shadow-md transition duration-200 ease-in-out hover:translate-y-0 hover:shadow-sm disabled:translate-y-0.5 disabled:opacity-70 disabled:shadow-none;
  }
}

.form-group-btn {
  @apply grid grid-rows-2 gap-2 pb-4 md:grid-cols-2 md:grid-rows-none;

  button {
    @apply -translate-y-0.5 rounded-lg p-2 font-semibold text-white shadow-md transition duration-200 ease-in-out first:bg-green-500 last:bg-red-600 hover:translate-y-0 hover:shadow-sm;
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
