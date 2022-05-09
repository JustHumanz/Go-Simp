<script setup>
import GroupPage from "../components/NewVtuber/GroupPage.vue"
import CreateGroup from "../components/NewVtuber/CreateGroup.vue"
</script>

<template>
  <div class="title">
    <span class="title__span">
      <font-awesome-icon icon="circle-plus" class="title__svg fa-fw" />
      Add new vtuber
    </span>
  </div>
  <div class="form">
    <transition name="close" @before-enter="test">
      <div v-if="step === 1" class="select-group">
        <group-page @group="getGroup" :groups="groups" />
      </div>
    </transition>
    <transition name="close" @before-enter="test">
      <div v-if="step === 2" class="create-group">
        <create-group @group="getGroup" @back="backAction" />
      </div>
    </transition>
    <transition name="close">
      <div v-if="step === 3"></div>
    </transition>
  </div>
</template>

<script>
import { library } from "@fortawesome/fontawesome-svg-core"
import { faCirclePlus } from "@fortawesome/free-solid-svg-icons"

library.add(faCirclePlus)

import axios from "axios"
import Config from "../config.json"

export default {
  data() {
    return {
      group: null,
      step: 1,
      groups: [],
    }
  },
  async mounted() {
    await this.getGroupData()
    // this.checkHeightDiv()
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
    },
    getGroup(group) {
      if (group.ID === -1) this.step = 2
      else this.step = 3
      this.group = group
    },
    async checkHeightDiv() {
      // await new Promise((resolve) => setTimeout(resolve, 500))
      // const getchild = document.querySelector(".form").children
      // console.log(getchild)
      // const childs = [...getchild]
      // // add var style with feiled --height-div
      // childs.forEach((child) => {
      //   const totalHeight = child.offsetHeight
      //   child.style.setProperty("--totalHeight", `${totalHeight}px`)
      // })
    },
    backAction() {
      this.step -= 1
    },
    test(e) {
      console.log(e)
    },
  },
}
</script>

<style lang="scss" scoped>
.title {
  @apply fixed top-16 z-[9] flex w-full select-none flex-wrap bg-blue-400 py-2 text-2xl font-semibold uppercase dark:bg-slate-500;

  &__span {
    @apply mx-auto w-[95%] text-white md:w-[75%] lg:w-[70%];
  }
  &__svg {
    @apply text-blue-200 dark:text-gray-200;
  }
}

.form {
  @apply mx-auto mt-28 w-[95%] overflow-x-hidden sm:w-[90%] md:w-[80%] lg:w-[75%];
}

.close {
  &-enter-to,
  &-leave-form {
    @apply scale-y-100;
  }
  &-enter-from,
  &-leave-to {
    @apply scale-y-0;
  }
  &-enter-active,
  &-leave-active {
    transition-property: transform height;
    @apply origin-top duration-300 ease-in-out;
  }
}
</style>
