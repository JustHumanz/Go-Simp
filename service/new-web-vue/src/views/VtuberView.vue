<template>
  <section class="vtuber-view">
<div class="photo-profile">
  <img v-if="member.Youtube" :src="member.Youtube.Avatar" alt="" class="photo-profile__img">
  <img v-else-if="member.BiliBili" :src="member.BiliBili.Avatar" alt="" class="photo-profile__img">
</div>
  </section>
</template>

<style lang="scss" scoped></style>

<script>
import axios from "axios"
import Config from "../config.json"

export default {
  data() {
    return {
      member: {},
    }
  },
  async mounted() {
    document.title = "Vtuber Details"
    await axios
      .get(Config.REST_API + "/members/" + this.$route.params.id)
      .then((response) => (this.member = response.data[0]))
console.log(this.member)
    document.title = this.member.EnName + " - Vtuber Details"
  },
}
</script>
