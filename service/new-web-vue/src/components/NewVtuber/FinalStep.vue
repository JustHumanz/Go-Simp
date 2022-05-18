<template>
  <h1 class="title-final">Completing in GitHub Issues</h1>
  <p class="description-final">
    All fields automatically convert to markdown for GitHub Issues.
  </p>

  <div class="links">
    <a :href="url" target="_blank" rel="noopener noreferrer">Finish</a>
    <button type="button" class="back" @click="$emit(`back`, true)">
      Back
    </button>
  </div>
</template>
<script>
import trim from "validator/lib/trim"
export default {
  props: {
    vtuber: {
      type: Array,
      required: true,
    },
    group: {
      type: Object,
    },
    newGroup: {
      type: Object,
    },
  },
  mounted() {
    // const title = encodeURIComponent(
    //   `Add to ${this.group.GroupName ?? this.newGroup.GroupName}`
    // )
    // const issue = encodeURIComponent(this.convert)
    // console.log(
    //   `https://github.com/JustHumanz/Go-Simp/issues/new?title=${title}&body=${issue}&labels=enhancement&assignees=JustHumanz`
    // )
  },
  computed: {
    convert() {
      // convert this.vtuber to markdown
      const vtubers = this.vtuber.map((vtuber) => {
        console.log(vtuber)
        return `\`\`\`
        Name = ${vtuber.nickname}
        EN_Name = "${vtuber.name_en}"
        JP_Name = "${vtuber.name_jp}"
        Region = "${vtuber.region}"
        Fanbase = "${vtuber.fanbase}"

        [Twitter]
        Twitter_Fanart = "${vtuber.twitter.fanart_hashtag}"
        Twitter_Lewd = "${vtuber.twitter.lewd_hashtag}"
        Twitter_Username = "${vtuber.twitter.username}"

        [Youtube]
        Yt_ID = "${vtuber.platform.youtube.channel_id}"

        [BiliBili]
        BiliBili_Fanart = "${vtuber.platform.bilibili.bili_fanart}"
        BiliBili_ID = "${vtuber.platform.bilibili.space_id}"
        BiliRoom_ID = "${vtuber.platform.bilibili.live_id}"

        [Twitch]
        Twitch_Username = "${vtuber.platform.twitch.username}"
        \`\`\``
      })

      let newGroup = ""

      if (this.newGroup) {
        const platformYt = !this.newGroup.youtube
          ? ""
          : this.newGroup.youtube.map((yt) => {
              return `
            [Youtube]
            Yt_ID = "${yt.channel_id}"
            `
            })

        const platformBili = !this.newGroup.bilibili
          ? ""
          : this.newGroup.bilibili.map((bili) => {
              return `
          [BiliBili]
          BiliBili_ID = "${bili.BiliBili_ID}"
          BiliRoom_ID = "${bili.BiliRoom_ID}"
          `
            })

        newGroup = `# New Group
          \`\`\`
      Group_Name = "${this.newGroup.GroupName}"
      Group_Icon = "${this.newGroup.IconUrl}"

      [Platforms]

      ${platformYt}

      ${platformBili}
      \`\`\``
      }

      return this.trims(`${newGroup}

      ${vtubers.join("\n\n")}
      `)
    },
    url() {
      const title = encodeURIComponent(
        `Add to ${this.group.GroupName ?? this.newGroup.GroupName}`
      )
      const issue = encodeURIComponent(this.convert)

      return `https://github.com/JustHumanz/Go-Simp/issues/new?title=${title}&body=${issue}&labels=enhancement&assignees=JustHumanz`
    },
  },
  methods: {
    trims(str) {
      // split string by new line
      const lines = str.split("\n")
      // trim
      const trimmed = lines.map((line) => {
        return trim(line)
      })

      return trimmed.join("\n")
    },
  },
}
</script>
<style lang="scss" scoped>
.title-final {
  @apply my-2 text-lg font-semibold;
}

.description-final {
  @apply my-2 text-gray-600;
}

.links {
  @apply grid grid-rows-2 gap-2 pb-4 md:w-3/6 md:grid-cols-2 md:grid-rows-none;

  a,
  button {
    @apply -translate-y-0.5 rounded-lg p-2 text-center font-semibold text-white shadow-md transition duration-200 ease-in-out first:bg-green-500 last:bg-red-600 hover:translate-y-0 hover:shadow-sm disabled:translate-y-0 disabled:opacity-70 disabled:shadow-none;
  }
}
</style>
