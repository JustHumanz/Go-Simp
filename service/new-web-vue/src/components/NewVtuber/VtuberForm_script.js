import Regions from "@/regions.json"
import trim from "validator/lib/trim"

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
    newVtubers: {
      type: Array,
      default: [],
    },
  },
  emits: ["back", "vtuber"],
  async mounted() {
    if (this.newVtubers.length > 0) this.assignNewVtuber(this.newVtubers)
    else this.addVtuber()
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
        const inputItem = vtuber.querySelectorAll(".vtuber__content-item")
        const oldInputItem = filteredvtubers[index].querySelectorAll(
          ".vtuber__content-item"
        )

        vtuber.classList = filteredvtubers[index].classList
        // move var --contentHeight
        const oldCalcHeight =
          filteredvtubers[index].lastElementChild.style.getPropertyValue(
            "--contentHeight"
          )
        vtuber.lastElementChild.style.setProperty(
          "--contentHeight",
          oldCalcHeight
        )

        inputItem.forEach((inp, i) => {
          // inp.value = oldInput[i].value
          inp.classList = oldInputItem[i].classList
          inp.querySelector(".error").innerHTML =
            oldInputItem[i].querySelector(".error").innerHTML
          inp.querySelector("input").value =
            oldInputItem[i].querySelector("input").value
        })
      })

      this.disabledDelBtn()
    },
    requestVtuber(e) {
      e.preventDefault()
      const result = []
      const vtubers = document.querySelectorAll(".vtuber-form")
      vtubers.forEach((vtuber) => {
        let vtuberData = {}

        vtuberData.nickname = trim(vtuber.querySelector("[name=name]").value)
        vtuberData.name_en = trim(vtuber.querySelector("[name=en-name]").value)
        vtuberData.name_jp = trim(vtuber.querySelector("[name=jp-name]").value)
        vtuberData.fanbase = trim(vtuber.querySelector("[name=fanbase]").value)
        vtuberData.region = Regions.find(
          (region) =>
            region.name === trim(vtuber.querySelector("[name=lang-code]").value)
        ).code

        vtuberData.platform = {
          youtube: {
            channel_id: trim(vtuber.querySelector("[name=youtube-id]").value),
          },
          twitch: {
            username: trim(
              vtuber.querySelector("[name=nickname-twitch]").value
            ),
          },
          bilibili: {
            space_id: trim(vtuber.querySelector("[name=space-id]").value),
            live_id: trim(vtuber.querySelector("[name=live-id]").value),
            bili_fanart: trim(vtuber.querySelector("[name=bili-art]").value),
          },
        }

        vtuberData.twitter = {
          username: trim(vtuber.querySelector("[name=twitter-username]").value),
          fanart_hashtag: trim(
            vtuber.querySelector("[name=fanart-hashtag]").value
          ),
          lewd_hashtag: trim(vtuber.querySelector("[name=lewd-hashtag]").value),
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

    async assignNewVtuber(vtubers) {
      if (!vtubers) return

      console.log(vtubers)

      for (const vtuber of vtubers) {
        this.vtubers.push({ id: this.vtubers.length, error: false })
      }

      await new Promise((resolve) => setTimeout(resolve, 60))
      const vtuberForms = document.querySelectorAll(".vtuber-form")

      vtuberForms.forEach((platform) => {
        platform.classList.remove("show")
      })

      let index = 0

      for (const form of vtuberForms) {
        const Inputs = form.querySelectorAll("input")
        const vtuberText = form.querySelector(".vtuber-link__text")

        vtuberText.innerText = vtubers[index].name_en
        Inputs[0].value = vtubers[index].nickname
        Inputs[1].value = vtubers[index].name_en
        Inputs[2].value = vtubers[index].name_jp
        Inputs[3].value = vtubers[index].fanbase
        Inputs[4].value = Regions.find(
          (region) => region.code === vtubers[index].region
        ).name
        Inputs[5].value = vtubers[index].platform.youtube?.channel_id
        Inputs[6].value = vtubers[index].platform.twitch?.username
        Inputs[7].value = vtubers[index].platform.bilibili?.space_id
        Inputs[8].value = vtubers[index].platform.bilibili?.live_id
        Inputs[9].value = vtubers[index].platform.bilibili?.bili_fanart
        Inputs[10].value = vtubers[index].twitter?.username
        Inputs[11].value = vtubers[index].twitter?.fanart_hashtag
        Inputs[12].value = vtubers[index].twitter?.lewd_hashtag

        index++
      }

      // get submit button
      const submitButton = document.querySelector(".submit")
      submitButton.disabled = this.vtubers.find((vtuber) => vtuber.error)
    },
  },
}
