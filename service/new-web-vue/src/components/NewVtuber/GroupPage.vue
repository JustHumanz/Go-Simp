<template>
  <h4 class="mt-2 text-lg font-semibold">Select your group...</h4>
  <ul class="group-list">
    <li class="group-list-item" v-for="group in groups">
      <a
        class="group-list-item__link"
        href="#"
        :data-id="group.ID"
        onclick="return false"
        @click="setGroup"
      >
        <font-awesome-icon
          class="fa-fw group-list-item__svg"
          icon="users"
          v-if="!group.ID"
        />
        <font-awesome-icon
          class="fa-fw group-list-item__svg"
          icon="user"
          v-else-if="group.ID && group.ID == 10"
        />
        <img
          draggable="false"
          v-else-if="group.ID && group.ID != 10"
          :src="group.GroupIcon"
          :alt="group.GroupName"
          class="group-list-item__img"
        />
        <span class="group-list-item__span">
          {{
            (
              group.GroupName.charAt(0).toUpperCase() + group.GroupName.slice(1)
            ).replace("_", " ")
          }}
        </span>
      </a>
    </li>
    <li class="group-list-item">
      <a
        href="#"
        class="group-list-item__link"
        onclick="return false"
        @click="setGroup"
      >
        <font-awesome-icon
          class="fa-fw group-list-item__svg fa-fw"
          icon="circle-plus"
        />
        <span class="group-list-item__span">Request new group</span>
      </a>
    </li>
  </ul>
</template>
<script>
import axios from "axios"
import Config from "../../config.json"

export default {
  data() {
    return {
      groups: [],
      toggle: false,
    }
  },
  async created() {
    const checkGroup = await this.getGroupData()

    if (!checkGroup) {
      // err here
    }

    console.log(this.groups)
  },
  methods: {
    async getGroupData() {
      if (this.groups.length > 0) return
      console.log("Fetching group data...")

      // this.cancelGroups = axios.CancelToken.source()

      const data_groups = await axios
        .get(Config.REST_API + "/v2/groups/", {
          // cancelToken: this.cancelGroups.token,
        })
        .then((response) => response.data)
        .catch((error) => {
          if (!axios.isCancel(error)) this.error_msg = error.message
        })

      if (this.error_msg) return false

      // sort group data from GroupName
      data_groups.sort((a, b) => {
        if (a.GroupName.toLowerCase() < b.GroupName.toLowerCase()) return -1
        if (a.GroupName.toLowerCase() > b.GroupName.toLowerCase()) return 1
        return 0
      })

      this.groups = data_groups
      console.log(`Total group: ${this.groups.length}`)
      await new Promise((resolve) => setTimeout(resolve, 500))
      this.checkHeightDiv()
      return true
    },
    async setGroup(e) {
      const id = e.target.closest(".group-list-item__link").dataset.id
      const group = this.groups.find((group) => group.ID == id)
      if (!group) this.$emit("group", { ID: -1 })
      else this.$emit("group", group)
    },
    checkHeightDiv() {
      // const getchild = document.querySelector(".form").children
      // console.log(getchild)
      // const childs = [...getchild]
      // // add var style with feiled --height-div
      // childs.forEach((child) => {
      //   const totalHeight = child.offsetHeight
      //   child.style.setProperty("--totalHeight", `${totalHeight}px`)
      // })
    },
  },
}
</script>
<style lang="scss" scoped>
.group {
  &-list {
    @apply grid gap-2 pt-2;
    grid-template-columns: repeat(auto-fit, minmax(11rem, 1fr));

    &-item {
      @apply overflow-hidden rounded-md bg-blue-400;

      &__link {
        @apply flex h-full items-center space-x-1 px-2 py-1 font-semibold text-white transition-all duration-200 ease-in-out hover:bg-black/10;
      }

      &__img {
        @apply w-6 object-contain drop-shadow-md;
      }

      &__svg {
        @apply w-6;
      }

      &__span {
        @apply text-sm font-semibold text-white;
      }
    }
  }
}
</style>
