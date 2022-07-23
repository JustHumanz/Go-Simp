<template>
  <h1 class="title-final">Completing in GitHub Issues</h1>
  <p class="description-final">
    All fields automatically convert to markdown for GitHub Issues.
  </p>

  <h3 class="title-preview">Preview</h3>
  <div class="my-contents" v-html="html"></div>

  <div class="links">
    <a :href="url" target="_blank" rel="noopener noreferrer">Finish</a>
    <button type="button" class="back" @click="$emit(`back`, true)">
      Back
    </button>
  </div>
</template>
<script>
import trim from "validator/lib/trim"
import { marked } from "marked"
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
  emits: ["back"],
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

        newGroup = `### New Group
          \`\`\`
      Group_Name = "${this.newGroup.GroupName}"
      Group_Icon = "${this.newGroup.IconUrl}"
${
  !this.newGroup.youtube && !this.newGroup.bilibili
    ? ""
    : `      [Platforms]

      ${platformYt}

      ${platformBili}
      `
}\`\`\``
      }

      return this.trims(`${newGroup}

### New Vtuber
      ${vtubers.join("\n\n")}
      `)
    },
    html() {
      const htmlData = marked(this.convert)
      console.log(htmlData)
      const getTag = /(<[^/][^>]+)(>)/gm
      const scopeId = this.$options.__scopeId

      return htmlData.replace(getTag, (match, p1, p2) => {
        return p1 + ` ${scopeId}=""` + p2
      })
    },
    url() {
      const title = encodeURIComponent(
        `Add to ${this.group?.GroupName ?? this.newGroup?.GroupName}`
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

.title-preview {
  @apply my-2 font-semibold;
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

.my-contents {
  @apply mb-4 space-y-2 overflow-y-auto rounded-md bg-slate-200 p-2 dark:bg-gray-500 md:h-72 md:w-3/6;

  h3 {
    @apply text-lg font-semibold;
  }

  pre {
    @apply rounded-md bg-slate-300 py-1 px-2 text-xs dark:bg-slate-600;
  }
}
</style>
