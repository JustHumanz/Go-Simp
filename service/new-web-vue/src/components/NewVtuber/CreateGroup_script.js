import trim from "validator/lib/trim"
import isUrl from "validator/lib/isURL"
export default {
  data() {
    return {
      platforms: [],
    }
  },
  props: {
    groups: {
      type: Array,
      default: [],
    },
    NewGroup: {
      type: Object,
      default: {},
    },
  },
  emits: ["group", "back"],
  async mounted() {
    await this.assignNewGroup(this.NewGroup)

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
      }

      activeElement = e.target
    })
  },
  methods: {
    sendGroup(e) {
      e.preventDefault()

      const form = e.target
      const platforms = document.querySelectorAll(".platform-group__content")

      const NewGroup = new Object()
      const GroupChannel = { youtube: [], bilibili: [] }

      NewGroup.GroupName = trim(form[0].value)
      NewGroup.IconUrl = trim(form[1].value)

      platforms.forEach((platform) => {
        const select = platform.querySelector("select")
        const input = platform.querySelectorAll("input")

        if (select.value === "youtube")
          GroupChannel.youtube.push({
            ChannelID: trim(input[0].value),
            Region: trim(input[1].value),
          })
        else if (select.value === "bilibili")
          GroupChannel.bilibili.push({
            BiliBili_ID: trim(input[0].value),
            BiliRoom_ID: trim(input[1].value),
            Region: trim(input[2].value),
          })
      })

      if (!GroupChannel.youtube.length) GroupChannel.youtube = null
      if (!GroupChannel.bilibili.length) GroupChannel.bilibili = null

      NewGroup.GroupChannel = GroupChannel
      this.$emit("group", NewGroup)
    },
    async addPlatform(e) {
      if (this.platforms.length == 9) e.target.disabled = true
      this.platforms.push({
        id: this.platforms.length,
        set: "youtube",
      })

      await new Promise((resolve) => setTimeout(resolve, 60))
      this.checkAllFilled()
    },
    async deletePlatform(id) {
      const addPlatformBtn = document.querySelector("[name=add-platform]")
      if (addPlatformBtn.disabled) addPlatformBtn.disabled = false

      const platforms = [...document.querySelector(".platforms").children]
      // remove confirm class from delete-platform
      platforms.forEach((platform) => {
        platform.querySelector(".delete-platform").classList.remove("confirm")
      })

      const filteredPlatforms = platforms.filter((platform, index) => {
        return index !== id
      })
      this.platforms.splice(id, 1)

      await new Promise((resolve) => setTimeout(resolve, 90))
      this.checkAllFilled()

      if (!filteredPlatforms.length) return

      this.platforms.map((platform, index) => {
        platform.id = index
      })

      filteredPlatforms.forEach((platform, index) => {
        const select = platform.querySelector("select")
        this.platforms[index].set = select.value
      })

      const newPlatforms = [...document.querySelector(".platforms").children]

      newPlatforms.forEach((platform, index) => {
        const input = platform.querySelectorAll("input")
        const oldInput = filteredPlatforms[index].querySelectorAll("input")

        platform.classList = filteredPlatforms[index].classList

        input.forEach((inp, i) => (inp.value = oldInput[i].value))
      })
    },
    checkFilled(element) {
      const notInput = element.tagName !== "INPUT"
      this.calculateHeight(element)

      if (notInput) return true
      return this.checkValidate(element)
    },
    checkAllFilled() {
      const inputs = document.querySelectorAll("input")
      let count = 0

      for (const input of inputs) {
        if (!this.checkValidate(input)) break
        count++
      }

      // get submit button
      const submitBtn = document.querySelector("[type=submit]")
      submitBtn.disabled = count < inputs.length
    },
    checkValidate(element) {
      const errorText = element.parentElement.querySelector(".error")
      const iconUrl = element.name === "group-icon"
      const value = trim(element.value)

      // toggle error class in parent element
      if (!iconUrl && !value) {
        errorText.innerText = "This field is required"
        return false
      }

      // when group name same
      if (element.name === "group-name") {
        const groupNameExist = this.groups.some(
          (group) =>
            group.GroupName.toLowerCase().replace("_", " ") ===
            value.toLowerCase()
        )

        if (groupNameExist) {
          errorText.innerText = "Group name already exist, please find new one"
          return false
        }
      }

      // when is url
      if (iconUrl) {
        if (value && !isUrl(value)) {
          errorText.innerText = "Please enter a valid URL (image)"
          return false
        }
      }

      // check is valid youtube id
      if (element.name === "youtube-id") {
        const isChannelId = value.match(/^UC[a-zA-Z0-9-_]{22}$/)

        if (!isChannelId) {
          errorText.innerText =
            "This is not a valid channel ID, find inside URL"
          return false
        }
      }

      // check is valid bilibili live id or space id
      if (element.name === "live-id") {
        const isBiliBiliLiveId = value.match(/^\d+$/)

        if (!isBiliBiliLiveId) {
          errorText.innerText = "This is not a valid live ID, find inside URL"
          return false
        }
      }

      if (element.name === "space-id") {
        const isBiliBiliSpaceId = value.match(/^\d+$/)

        if (!isBiliBiliSpaceId) {
          errorText.innerText = "This is not a valid space ID, find inside URL"
          return false
        }
      }

      // check region/language code
      if (element.name === "lang-code") {
        const isRegionCode = value.match(/^[a-zA-Z]{2}$/)

        if (!isRegionCode) {
          errorText.innerText = "Please enter a valid region/language code"
          return false
        }
      }
      return true
    },
    async calculateHeight(element) {
      const outsideDiv = element.parentElement.parentElement

      if (!outsideDiv?.classList.contains("platform-group__content")) return
      await new Promise((resolve) => setTimeout(resolve, 60))

      const calculatedHeight = [...outsideDiv.children].reduce((t, c) => {
        // get margin and padding
        const margin =
          parseInt(getComputedStyle(c).marginTop) +
          parseInt(getComputedStyle(c).marginBottom)
        const padding =
          parseInt(getComputedStyle(c).paddingTop) +
          parseInt(getComputedStyle(c).paddingBottom)

        // get total height
        return t + c.offsetHeight + margin + padding
      }, 0)

      outsideDiv.style.setProperty("--contentHeight", `${calculatedHeight}px`)
    },
    async assignNewGroup(group) {
      console.log(group)

      if (!group) return

      const youtubeChannels = group.GroupChannel.youtube
      const bilibiliChannels = group.GroupChannel.bilibili

      if (youtubeChannels) {
        for (const youtube of youtubeChannels) {
          this.platforms.push({
            id: this.platforms.length,
            set: "youtube",
          })
        }
      }

      if (bilibiliChannels) {
        for (const bilibili of bilibiliChannels) {
          this.platforms.push({
            id: this.platforms.length,
            set: "bilibili",
          })
        }
      }

      await new Promise((resolve) => setTimeout(resolve, 60))
      const platformGroup = document.querySelectorAll(".platform-group")

      platformGroup.forEach((platform) => {
        platform.classList.remove("show")
      })

      // fill all inputs

      const nameInput = document.querySelector("[name=group-name]")
      const iconInput = document.querySelector("[name=group-icon]")

      nameInput.value = group.GroupName
      iconInput.value = "" ?? group.GroupIcon

      let indexyt = 0
      let indexbili = 0

      platformGroup.forEach((platform) => {
        const platformSelect = platform.querySelector("select")
        const platformInputs = platform.querySelectorAll("input")

        platformInputs[0].value =
          platformSelect.value === "youtube"
            ? youtubeChannels[indexyt].ChannelID
            : bilibiliChannels[indexbili].BiliBili_ID

        platformInputs[1].value =
          platformSelect.value === "youtube"
            ? youtubeChannels[indexyt].Region
            : bilibiliChannels[indexbili].BiliRoom_ID

        if (platformSelect.value === "bilibili")
          platformInputs[2] = bilibiliChannels[indexbili].Region

        if (platformSelect.value === "youtube") indexyt++
        else indexbili++
      })

      await new Promise((resolve) => setTimeout(resolve, 60))
      this.checkAllFilled()
    },
  },
}
