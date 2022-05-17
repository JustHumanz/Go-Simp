export default {
  data() {
    return {
      vtubers: [],
    }
  },
  props: {
    group: {
      type: Object,
      required: true,
    },
    nicknames: {
      type: Array,
      default: [],
    },
  },
  emits: ["back", "vtuber"],
  async mounted() {
    this.addVtuber()
  },
  methods: {
    async addVtuber(e = null) {
      if (e === null) e = document.querySelector("[name=add-vtuber]")
      if (this.vtubers.length == 14) e.target.disabled = true

      const vtuberForms = document.querySelectorAll(".vtuber-form")
      for (const vtuberForm of vtuberForms) {
        vtuberForm.classList.remove("show")
      }

      this.vtubers.push({ id: this.vtubers.length, error: true })

      await new Promise((resolve) => setTimeout(resolve, 60))
      this.disabledDelBtn()
    },
    checkError(data = null) {
      if (data) {
        this.vtubers.map((vtuber) => {
          if (vtuber.id == data.id) vtuber.error = data.error
        })
      }

      // get submit button
      const submitButton = document.querySelector(".submit")
      submitButton.disabled = this.vtubers.find((vtuber) => vtuber.error)
    },
    async deletePlatform(id) {
      const addVtuberBtn = document.querySelector("[name=add-vtuber]")
      if (addVtuberBtn.disabled) addVtuberBtn.disabled = false

      const vtubers = [...document.querySelector(".vtubers").children]
      // remove confirm class from delete-platform
      vtubers.forEach((vtuber) => {
        vtuber.querySelector(".delete-vtuber").classList.remove("confirm")
      })

      const filteredvtubers = vtubers.filter((platform, index) => {
        return index !== id
      })
      this.vtubers.splice(id, 1)

      await new Promise((resolve) => setTimeout(resolve, 90))

      if (!filteredvtubers.length) return

      this.vtubers.map((vtuber, index) => {
        vtuber.id = index
      })

      const newvtubers = [...document.querySelector(".vtubers").children]

      newvtubers.forEach((vtuber, index) => {
        const input = vtuber.querySelectorAll("input")
        const oldInput = filteredvtubers[index].querySelectorAll("input")

        vtuber.classList = filteredvtubers[index].classList

        input.forEach((inp, i) => (inp.value = oldInput[i].value))
      })

      this.disabledDelBtn()
    },
    requestVtuber(e) {
      e.preventDefault()
      const result = []
      const vtubers = document.querySelectorAll(".vtuber-form")
      vtubers.forEach((vtuber) => {
        let vtuberData = {}

        vtuberData.nickname = vtuber.querySelector("[name=name]").value
        vtuberData.name_en = vtuber.querySelector("[name=en-name]").value
        vtuberData.name_jp = vtuber.querySelector("[name=jp-name]").value
        vtuberData.fanbase = vtuber.querySelector("[name=fanbase]").value
        vtuberData.region = vtuber.querySelector("[name=lang-code]").value

        vtuberData.platform = {
          youtube: {
            channel_id: vtuber.querySelector("[name=youtube-id]").value,
          },
          twitch: {
            username: vtuber.querySelector("[name=nickname-twitch]").value,
          },
          bilibili: {
            space_id: vtuber.querySelector("[name=space-id]").value,
            live_id: vtuber.querySelector("[name=live-id]").value,
          },
        }

        vtuberData.twitter = {
          username: vtuber.querySelector("[name=twitter-username]").value,
          fanart_hashtag: vtuber.querySelector("[name=fanart-hashtag]").value,
          lewd_hashtag: vtuber.querySelector("[name=lewd-hashtag]").value,
        }

        result.push(vtuberData)
      })

      this.$emit("vtuber", result)
    },

    disabledDelBtn() {
      // get all delete buttons
      const deleteButtons = document.querySelectorAll(".delete-vtuber")
      for (const deleteButton of deleteButtons) {
        deleteButton.classList.toggle("one", this.vtubers.length == 1)
      }
    },
  },
}
