<script setup>
import GroupPage from "../components/NewVtuber/GroupPage.vue"
import CreateGroup from "../components/NewVtuber/CreateGroup.vue"
import VtuberForm from "../components/NewVtuber/VtuberForm.vue"
import FinalStep from "../components/NewVtuber/FinalStep.vue"
</script>

<template>
  <div class="title">
    <span class="title__span">
      <font-awesome-icon icon="circle-plus" class="title__svg fa-fw" />
      Add new vtuber
    </span>
  </div>
  <div class="form">
    <transition-group
      leave-from-class="slide-center"
      :leave-to-class="step > oldStep ? 'slide-left' : 'slide-right'"
      leave-active-class="slide-active"
      :enter-from-class="oldStep < step ? 'slide-right' : 'slide-left'"
      enter-to-class="slide-center"
      enter-active-class="slide-active"
    >
      <div v-if="step === 1" class="select-group" ref="container">
        <group-page @group="getGroup" :groups="groups" />
      </div>
      <div v-if="step === 2" class="create-group" ref="container">
        <create-group
          @group="newGroupfunction"
          @back="backAction"
          :groups="groups"
          :NewGroup="newGroup"
        />
      </div>
      <div v-if="step === 3">
        <vtuber-form
          :group="group ? group : newGroup"
          :nicknames="nicknames"
          :newVtubers="newVtuber"
          @back="backAction"
          @vtuber="newVtuberFunction"
        />
      </div>
      <div v-if="step === 4">
        <final-step
          :group="group"
          :newGroup="newGroup"
          :vtuber="newVtuber"
          @back="backAction"
        />
      </div>
    </transition-group>
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
      newGroup: null,
      newVtuber: [],
      step: 1,
      oldStep: 1,
      groups: [],
      nicknames: [],
    }
  },
  async mounted() {
    document.title = "Add New Vtuber - Vtbot"
    await this.getGroupData()
    this.getNickname()

    this.$watch(
      () => this.step,
      (newValue, oldValue) => (this.oldStep = oldValue)
    )

    this.$watch(
      () => this.group,
      () => {
        if (this.newVtuber.length > 0) {
          this.newVtuber = []
        }

        if (this.group) {
          this.newGroup = null
        }
      }
    )
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
      if (!group) this.step = 2
      else {
        this.newGroup = null
        this.step = 3
      }
      this.group = group
    },
    newGroupfunction(group) {
      this.newGroup = group
      this.step = 3
    },
    newVtuberFunction(vtuber) {
      this.newVtuber = vtuber
      this.step = 4
    },
    backAction() {
      if (this.group && this.step == 3) return (this.step = 1)
      // if (this.step === 2 && this.newGroup) this.newGroup = null
      this.step -= 1
    },
    async getNickname() {
      const vtuber_data = await axios
        .get(Config.REST_API + "/v2/members/", {})
        .then((response) => {
          return response.data.map((vtuber) => vtuber.NickName)
        })
        .catch((error) => {
          if (!axios.isCancel(error)) this.error_status = error.response.status
        })
      if (!vtuber_data) return

      this.nicknames = vtuber_data
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

.slide {
  &-left {
    @apply -translate-x-[150%];
  }

  &-right {
    @apply translate-x-[150%];
  }

  &-center {
    @apply translate-x-0;
  }

  &-active {
    transition-property: transform height;
    @apply absolute w-[95vw] overflow-hidden duration-300 ease-in-out sm:w-[90vw] md:w-[80vw] lg:w-[75vw];
  }
}
</style>
