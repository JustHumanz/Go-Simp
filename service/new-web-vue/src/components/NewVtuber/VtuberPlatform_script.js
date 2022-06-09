import trim from "validator/lib/trim"
import Regions from "@/regions.json"

export default {
  data() {
    return {
      vtuberName: "",
      toggleLang: null,
      searchLang: "",
    }
  },
  props: {
    id: {
      type: Number,
      required: true,
    },
    nicknames: {
      type: Array,
      default: [],
    },
  },
  emits: ["error", "delete"],
  async mounted() {
    const vtuberForm = this.$refs.form

    vtuberForm.querySelectorAll(".vtuber__content-item").forEach((item) => {
      item.querySelector("input").value = ""
      item.classList.remove("has-error")
      item.querySelector(".error").innerHTML = ""
    })

    await this.checkHeight()

    vtuberForm.addEventListener("input", async (e) => {
      e.target.parentElement.classList.toggle(
        "has-error",
        !(await this.checkFilled(e.target, true))
      )

      await new Promise((resolve) => setTimeout(resolve, 60))
      await this.checkAllFilled(e.target)
      this.checkHeight()
    })

    let activeElement = null

    document.body.addEventListener("click", async (e) => {
      if (
        activeElement &&
        e.target !== activeElement &&
        activeElement.tagName === "INPUT" &&
        vtuberForm === activeElement?.closest(".vtuber-form")
      ) {
        activeElement.parentElement.classList.toggle(
          "has-error",
          !(await this.checkFilled(activeElement))
        )

        await new Promise((resolve) => setTimeout(resolve, 60))
        await this.checkAllFilled(activeElement)
        this.checkHeight()
      }

      activeElement = e.target
    })

    document.body.addEventListener("click", (e) => {
      if (!e.target.closest(".delete-vtuber")) {
        document.querySelectorAll(".delete-vtuber").forEach((vtuber) => {
          vtuber.classList.remove("confirm")
        })
      }
    })
  },
  computed: {
    getVtuberName() {
      return this.vtuberName ? this.vtuberName : "New Vtuber"
    },

    regions() {
      return Regions.filter(
        (region) =>
          region.code.toLowerCase().includes(this.searchLang.toLowerCase()) ||
          region.name.toLowerCase().includes(this.searchLang.toLowerCase())
      )
    },
  },

  methods: {
    async checkFilled(element, isInput = false) {
      const notInput = element.tagName !== "INPUT"
      // this.calculateHeight(element)
      if (notInput || isInput) return true
      return await this.checkValidate(element, isInput)
    },
    async checkAllFilled(e) {
      const form = this.$refs.form
      if (!form) return

      const inputs = form.querySelectorAll("input")
      let count = 0
      for (const input of inputs) {
        const validate = await this.checkValidate(input, true)
        if (!validate) break
        count++
      }

      const youtubeId = e?.name === "youtube-id"
      const twitchName = e?.name === "nickname-twitch"
      const biliSpaceId = e?.name === "space-id"
      const biliLiveId = e?.name === "live-id"

      const platforms =
        youtubeId || twitchName || biliSpaceId || biliLiveId ? e.name : null

      const isFilled = this.CheckPlatform(platforms)

      // console.log(count, inputs.length, isFilled)

      const error = count < inputs.length || !isFilled
      // console.log(error)
      this.$refs.form.classList.toggle("errors", error)
      this.$emit("error", { id: this.id, error })
    },
    async checkValidate(e, isInput = false) {
      const name = e.name === "name"
      const ENname = e.name === "en-name"
      const JPname = e.name === "jp-name"
      const region = e.name === "lang-code"
      const youtubeId = e.name === "youtube-id"
      const twitchName = e.name === "nickname-twitch"
      const biliSpaceId = e.name === "space-id"
      const biliLiveId = e.name === "live-id"
      const biliFanArt = e.name === "bili-art"
      const twitterUser = e.name === "twitter-username"
      const fanartHash = e.name === "fanart-hashtag"
      const lewdHash = e.name === "lewd-hashtag"

      const errorText = e.parentElement.querySelector(".error")
      const value = trim(e.value)

      // check when name, ENname, and region is empty
      if ((name || ENname || region) && !value) {
        errorText.innerText = "Please enter a value"
        return false
      }

      // check when nickname invalid
      const isNotNickname = value.match(/[^a-zA-Z0-9-_]/g)
      if (name && isNotNickname) {
        errorText.innerText = "Please enter a valid nickname"
        return false
      }

      // check nickname exist
      if (name) {
        await new Promise((resolve) => {
          const waitNicknames = setInterval(() => {
            if (this.nicknames.length > 0) {
              clearInterval(waitNicknames)

              resolve()
            }
          }, 80)
        })

        const getNickname = this.nicknames.find(
          (nickname) => nickname.toLowerCase() === value.toLowerCase()
        )

        if (getNickname) {
          errorText.innerText = "This nickname is already taken"
          return false
        }
      }

      // when jpname is not japanese character
      const isJapaneseChar = value.match(
        /[一-龠]+|[ぁ-ゔ]+|[ァ-ヴー]+|[々〆〤・]+/gu
      )
      if (JPname && value && !isJapaneseChar) {
        errorText.innerText = "Please enter a valid Japanese name"
        return false
      }

      //check is region code
      const isRegion = Regions.find((region) => region.name === value)

      if (region && !isRegion) {
        if (!isInput) errorText.innerText = "Please enter a valid region"
        return false
      }

      const isytId = value.match(/^UC[a-zA-Z0-9-_]{22}$/)
      const isTwitchName = value.match(/^[a-zA-Z0-9_]{1,20}$/)
      const isBiliBiliId = value.match(/^\d+$/)
      //no space char
      const noSpace = value.match(/\s/)

      // when youtubeId is invalid
      if (youtubeId && value && !isytId) {
        errorText.innerText = "Please enter a valid YouTube ID"
        return false
      }

      // when twitchName is invalid
      if (twitchName && value && !isTwitchName) {
        errorText.innerText = "Please enter a valid Twitch name"
        return false
      }

      // when biliSpaceId is invalid
      if (biliSpaceId && value && !isBiliBiliId) {
        errorText.innerText = "Please enter a valid BiliBili Space ID"
        return false
      }

      // when biliLiveId is invalid
      if (biliLiveId && value && !isBiliBiliId) {
        errorText.innerText = "Please enter a valid BiliBili Live ID"
        return false
      }

      // when biliFanArt is invalid
      if (biliFanArt && value && noSpace) {
        errorText.innerText = "Please enter a without space"
        return false
      }

      // check twitterUser is invalid
      const isTwitterUser = value.match(/^(@)?[a-zA-Z0-9_]*$/)
      if (twitterUser && value && !isTwitterUser) {
        errorText.innerText = "Please enter a valid Twitter username"
        return false
      }

      // check HashtagTwitter is invalid
      const isHashtagTwitter = value.match(
        /^(#)?[一-龠ぁ-ゔァ-ヴーa-zA-Z0-9々〆〤_]*$/
      )

      if ((fanartHash || lewdHash) && value && !isHashtagTwitter) {
        errorText.innerText = "Please enter a valid Twitter hashtag"
        return false
      }

      return true
    },
    CheckPlatform(name = null) {
      const youtube = this.$refs.ytid
      const twitch = this.$refs.twitchname
      const biliSpace = this.$refs.biliid
      const biliLive = this.$refs.liveid
      const parentYt = youtube.parentElement
      const parentTwitch = twitch.parentElement
      const parentBili = biliSpace.parentElement
      const parentBiliLive = biliLive.parentElement
      const errTextYt = parentYt.querySelector(".error")
      const errTextTwitch = parentTwitch.querySelector(".error")
      const errTextBili = parentBili.querySelector(".error")
      const errTextBiliLive = parentBiliLive.querySelector(".error")
      const isEmpty =
        !youtube.value && !twitch.value && (!biliSpace.value || !biliLive.value)
      if (isEmpty) {
        errTextYt.innerText = "Please at least add one platform"
        errTextTwitch.innerText = "Please at least add one platform"
        errTextBili.innerText = "Please at least add one platform"
        errTextBiliLive.innerText = "Please at least add one platform"
        if (name !== null) {
          parentYt.classList.add("has-error")
          parentTwitch.classList.add("has-error")
          parentBili.classList.add("has-error")
          parentBiliLive.classList.add("has-error")
        }
        return false
      } else {
        if (name !== "youtube-id") parentYt.classList.remove("has-error")
        if (name !== "nickname-twitch")
          parentTwitch.classList.remove("has-error")
        if (name !== "space-id") parentBili.classList.remove("has-error")
        if (name !== "live-id") parentBiliLive.classList.remove("has-error")
        return true
      }
    },
    async checkHeight() {
      if (!this.$refs.content) return
      await new Promise((resolve) => setTimeout(resolve, 60))
      const calculatedHeight = [...this.$refs.content.children].reduce(
        (t, c) => {
          // get margin and padding
          const margin =
            parseInt(getComputedStyle(c).marginTop) +
            parseInt(getComputedStyle(c).marginBottom)
          const padding =
            parseInt(getComputedStyle(c).paddingTop) +
            parseInt(getComputedStyle(c).paddingBottom)

          // get total height
          return t + c.offsetHeight + margin + padding
        },
        0
      )

      this.$refs.content.style.setProperty(
        "--contentHeight",
        `${calculatedHeight}px`
      )
    },
    deletePlatform(e) {
      const vtuberForms = document.querySelectorAll(".vtuber-form")
      if (vtuberForms.length < 2) return

      document.querySelectorAll(".delete-vtuber").forEach((platform) => {
        if (platform !== e.target.closest(".delete-vtuber"))
          platform.classList.remove("confirm")
      })
      if (e.target.closest(".delete-vtuber").classList.contains("confirm"))
        this.$emit("delete", this.id)
      else e.target.closest(".delete-vtuber").classList.add("confirm")
    },
    toggleContent(e) {
      if (e.target.closest(".delete-vtuber")) return
      const group = e.target.closest(".vtuber-form")

      const groups = [...document.querySelectorAll(".vtuber-form")].filter(
        (c) => c !== group
      )

      groups.forEach((c) => c.classList.remove("show"))

      group.classList.toggle("show")
    },

    openLang(e) {
      const input = e.target

      if (input.value && this.searchLang !== "")
        input.setSelectionRange(0, input.value.length)

      this.toggleLang = this.id
      this.searchLang = ""
    },

    setReg(e) {
      // get value radio
      const selectedRegion =
        e.target.closest(".region__label")?.previousElementSibling.value

      // set value to input
      const legionInput = e.target
        .closest(".vtuber__content-item")
        .querySelector("input[name='lang-code']")
      legionInput.value = Regions.find((r) => r.code === selectedRegion)?.name
      this.searchLang = Regions.find((r) => r.code === selectedRegion)?.name
      this.checkAllFilled(e.target)
    },
    findReg(e) {
      e.target.classList.toggle(
        "region-selected",
        !!this.regions.find((r) => r.name === e.target.value)
      )
      this.searchLang = e.target.value
    },
  },
}
