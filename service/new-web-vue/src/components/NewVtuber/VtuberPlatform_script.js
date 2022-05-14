import trim from "validator/lib/trim"

export default {
  data() {
    return {
      vtuberName: "",
    }
  },
  emits: ["error"],
  async mounted() {
    await new Promise((resolve) => setTimeout(resolve, 60))
    document.body.addEventListener("input", (e) => {
      e.target.parentElement.classList.toggle(
        "has-error",
        !this.checkFilled(e.target)
      )

      this.checkAllFilled()
    })

    let activeElement = null

    document.body.addEventListener("click", (e) => {
      if (activeElement && e.target !== activeElement) {
        activeElement.parentElement.classList.toggle(
          "has-error",
          !this.checkFilled(activeElement)
        )

        this.checkAllFilled()
      }

      activeElement = e.target
    })
  },
  computed: {
    getVtuberName() {
      return this.vtuberName ? this.vtuberName : "New Vtuber"
    },
  },

  methods: {
    checkFilled(element) {
      const notInput = element.tagName !== "INPUT"
      // this.calculateHeight(element)

      if (notInput) return true
      return this.checkValidate(element)
    },
    checkAllFilled() {
      const inputs = this.$refs.form.querySelectorAll("input")
      let count = 0

      for (const input of inputs) {
        if (!this.checkValidate(input, false)) break

        count++
      }

      console.log(count, inputs.length)

      const error = count < inputs.length
      console.log(error)

      this.$refs.form.classList.toggle("errors", error)

      this.$emit("error", error)
    },
    checkValidate(e, checked = true) {
      const name = e.name === "name"
      const ENname = e.name === "en-name"
      const JPname = e.name === "jp-name"
      // const fanbase = e.name === "fanbase"
      const region = e.name === "lang-code"
      const youtubeId = e.name === "youtube-id"
      const twitchName = e.name === "nickname-twitch"
      const biliSpaceId = e.name === "space-id"
      const biliLiveId = e.name === "live-id"
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

      // when jpname is not japanese character
      const isJapaneseChar = value.match(
        /[一-龠]+|[ぁ-ゔ]+|[ァ-ヴー]+|[々〆〤・]+/gu
      )
      if (JPname && value && !isJapaneseChar) {
        errorText.innerText = "Please enter a valid Japanese name"
        return false
      }

      //check is region code
      const isRegionCode = value.match(/^[a-zA-Z]{2}$/)
      if (region && !isRegionCode) {
        errorText.innerText = "Please enter a valid region/language code"
        return false
      }

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

        if (!youtubeId && checked) parentYt.classList.add("has-error")
        if (!twitchName && checked) parentTwitch.classList.add("has-error")
        if (!biliSpaceId && checked) parentBili.classList.add("has-error")
        if (!biliLiveId && checked) parentBiliLive.classList.add("has-error")

        if (youtubeId || twitchName || biliSpaceId || biliLiveId) return false
      } else {
        if (!youtubeId && checked) parentYt.classList.remove("has-error")
        if (!twitchName && checked) parentTwitch.classList.remove("has-error")
        if (!biliSpaceId && checked) parentBili.classList.remove("has-error")
        if (!biliLiveId && checked) parentBiliLive.classList.remove("has-error")
      }

      const isytId = value.match(/^UC[a-zA-Z0-9-_]{22}$/)
      const isTwitchName = value.match(/^[a-zA-Z0-9_]{1,20}$/)
      const isBiliBiliId = value.match(/^\d+$/)

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

      // check twitterUser is invalid
      const isTwitterUser = value.match(/^(@)?[a-zA-Z0-9_]$/)
      if (twitterUser && value && !isTwitterUser) {
        errorText.innerText = "Please enter a valid Twitter username"
        return false
      }

      // check HashtagTwitter is invalid
      const isHashtagTwitter = value.match(
        /^(#)?[一-龠ぁ-ゔァ-ヴーa-zA-Z0-9々〆〤_]$/
      )

      if ((fanartHash || lewdHash) && value && !isHashtagTwitter) {
        errorText.innerText = "Please enter a valid Twitter hashtag"
        return false
      }

      return true
    },
  },
}
