<script setup>
import VtuberCard from "./VtuberCard.vue"
</script>

<template>
  <section class="vtuber-list">
    <VtuberCard
      v-for="vtuber in limitedVtubers"
      :key="vtuber.ID"
      :vtuber="vtuber"
    />
  </section>
</template>

<script>
export default {
  props: {
    vtubers: {
      type: Array,
    },
  },
  data() {
    return {
      limitedVtubers: [],
    }
  },
  created() {
    this.$watch(
      () => this.$route,
      () => {
        this.limitedVtubers = this.sortingVtubers.slice(0, 25)
      },
      { immediate: true }
    )
  },
  computed: {
    sortingVtubers() {
      // get query.sort when exist
      const sort = this.$route.query.sort ? this.$route.query.sort : ""
      let vtuber_data = this.vtubers

      switch (sort) {
        case "":
          console.log("Sort by Alphabet")
          // Sorting vtuber by EnName ASC and lowercase
          vtuber_data = vtuber_data.sort((a, b) => {
            if (a.EnName.toLowerCase() < b.EnName.toLowerCase()) return -1
            if (a.EnName.toLowerCase() > b.EnName.toLowerCase()) return 1
            return 0
          })
          break
        case "-name":
          console.log("Sort by Reverse Alphabet")
          // Sorting vtuber by EnName DESC and lowercase
          vtuber_data = vtuber_data.sort((a, b) => {
            if (a.EnName.toLowerCase() > b.EnName.toLowerCase()) return -1
            if (a.EnName.toLowerCase() < b.EnName.toLowerCase()) return 1
            return 0
          })
          break
        case "yt":
          console.log("Sort by Most Youtube Subscriber")
          // Sorting vtuber by Youtube.Subscriber DESC when exist
          vtuber_data = vtuber_data.sort((a, b) =>{
            if (a.Youtube?.Subscriber < b.Youtube?.Subscriber) return 1
            if (a.Youtube?.Subscriber > b.Youtube?.Subscriber) return -1
            return 0
          })
          break
        case "-yt":
          console.log("Sort by Least Youtube Subscriber")
          // Sorting vtuber by Youtube.Subscriber ASC when exist
          vtuber_data = vtuber_data.sort((a, b) => {
            if (a.Youtube?.Subscriber > b.Youtube?.Subscriber) return -1
            if (a.Youtube?.Subscriber < b.Youtube?.Subscriber) return 1
            return 0
          })
          break
        case "tw":
          console.log("Sort by Most Twitch Followers")
          //Sorting vtuber by Twitch.Followers DESC when exist
          vtuber_data = vtuber_data.sort((a, b) => {
            if (a.Twitch?.Followers < b.Twitch?.Followers) return 1
            if (a.Twitch?.Followers > b.Twitch?.Followers) return -1
            return 0
          })
          break
        case "-tw":
          console.log("Sort by Least Twitch Followers")
          //Sorting vtuber by Twitch.Followers ASC when exist
          vtuber_data = vtuber_data.sort((a, b) => {
            if (a.Twitch?.Followers > b.Twitch?.Followers) return -1
            if (a.Twitch?.Followers < b.Twitch?.Followers) return 1
            return 0
          })
          break
        case "bl":
          console.log("Sort by Most Bilibili Followers")
          //Sorting vtuber by BiliBili.Followers DESC when exist
          vtuber_data = vtuber_data.sort((a, b) => {
            if (a.BiliBili?.Followers < b.BiliBili?.Followers) return 1
            if (a.BiliBili?.Followers > b.BiliBili?.Followers) return -1
            return 0
          })
          break
        case "-bl":
          console.log("Sort by Least Bilibili Followers")
          //Sorting vtuber by BiliBili.Followers ASC when exist
          vtuber_data = vtuber_data.sort((a, b) => {
            if (a.BiliBili?.Followers > b.BiliBili?.Followers) return -1
            if (a.BiliBili?.Followers < b.BiliBili?.Followers) return 1
            return 0
          })
          break
        case "twr":
          console.log("Sort by Most Twitter Followers")
          //Sorting vtuber by Twitter.Followers DESC when exist
          vtuber_data = vtuber_data.sort((a, b) => {
            if (a.Twitter?.Followers < b.Twitter?.Followers) return 1
            if (a.Twitter?.Followers > b.Twitter?.Followers) return -1
            return 0
          })
          break
        case "-twr":
          console.log("Sort by Least Twitter Followers")
          //Sorting vtuber by Twitter.Followers ASC when exist
          vtuber_data = vtuber_data.sort((a, b) => {
            if (a.Twitter?.Followers > b.Twitter?.Followers) return -1
            if (a.Twitter?.Followers < b.Twitter?.Followers) return 1
            return 0
          })
          break
      }

      // Sorting vtuber when IsLive.BiliBili is not null (Object), IsLive.Twitch is not null (Object), and IsLive.Youtube is not null (Object)
      return vtuber_data.sort((a, b) => {
        if (!a.IsLive.Youtube && b.IsLive.Youtube) return 1
        if (a.IsLive.Youtube && !b.IsLive.Youtube) return -1
        if (!a.IsLive.Twitch && b.IsLive.Twitch) return 1
        if (a.IsLive.Twitch && !b.IsLive.Twitch) return -1
        if (!a.IsLive.BiliBili && b.IsLive.BiliBili) return 1
        if (a.IsLive.BiliBili && !b.IsLive.BiliBili) return -1
        return 0
      })
    },
  },
}
</script>

<style lang="scss" scoped>
.vtuber-list {
  @apply pt-24 pb-4 xs:pt-14 w-full md:w-[80%] lg:w-[75%] mx-auto grid gap-[1.5rem];
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
}
</style>
