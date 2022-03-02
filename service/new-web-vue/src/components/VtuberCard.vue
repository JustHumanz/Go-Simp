<template>
  <div class="card-vtuber">
    {{ vtuber.EnName }}
    {{ vtuber.Youtube?.Subscriber || 0 }}
  </div>
</template>

<script>
export default {
  props: {
    vtuber: {
      type: Object,
    },
  },
  created() {
    // console.log("Show Vtuber: " + this.vtuber.EnName)
  },
}
</script>

<style lang="scss" scoped></style>
