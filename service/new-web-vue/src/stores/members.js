import { defineStore } from "pinia"
import { ref, computed, compile, toRaw } from "vue"
import axios from "axios"
import Config from "../config.json"
import regionConfig from "../regions.json"
import parse from "url-parse"

export const useMemberStore = defineStore("members", () => {
  const members = ref({
    error: false,
    query: "",
    status: null,
    group: -1,
    data: [],
    filteredData: [],
    searchedData: [],
  })

  const menuFilter = ref({
    region: [],
    platform: [],
    live: [],
    inactive: false,
  })
  const sortMenu = ref({ platform: [], twitter: false })

  const filter = ref({
    region: null,
    platform: null,
    live: null,
    live_only: false,
    inactive: false,
  })
  const sorting = ref({
    type: "name",
    order: "asc",
    live: true,
    inactive: true,
  })

  const fetchMembers = async (id = null) => {
    let err = false
    members.value.status = null
    members.value.error = false

    members.value.data = []
    members.value.filteredData = []

    members.value.group = id

    const params = id ? { groupid: id } : {}

    const vtuber_data = await axios
      .get(Config.REST_API + "/v2/members/", {
        params,
      })
      .then((response) => response.data)
      .catch((error) => {
        err = true
        members.value.error = true
        members.value.status = error.response.status
      })

    if (err) return false

    vtuber_data.forEach((vtuber) => {
      regionConfig.forEach((region) => {
        if (region.code === vtuber.Region) {
          vtuber.Regions = region
        }
      })
    })

    console.log(`Total member: ${vtuber_data.length}`)

    vtuber_data.sort(
      // sort by name
      ({ EnName: nameA }, { EnName: nameB }) =>
        nameA.toLowerCase() < nameB.toLowerCase() ? -1 : 1
    )

    let newRegion = []
    let newPlatform = []
    let newLive = []
    let newInac = false

    for (const vt of vtuber_data) {
      if (!newRegion.includes(vt.Region)) newRegion.push(vt.Region)

      if (!newPlatform.includes("youtube") && vt.Youtube)
        newPlatform.push("youtube")
      if (!newPlatform.includes("twitch") && vt.Twitch)
        newPlatform.push("twitch")
      if (!newPlatform.includes("bilibili") && vt.BiliBili)
        newPlatform.push("bilibili")

      if (!newLive.includes("youtube") && vt.IsLive.Youtube)
        newLive.push("youtube")
      if (!newLive.includes("twitch") && vt.IsLive.Twitch)
        newLive.push("twitch")
      if (!newLive.includes("bilibili") && vt.IsLive.BiliBili)
        newLive.push("bilibili")

      if (vt.Status === "Inactive") newInac = true
    }

    menuFilter.value = {
      region: newRegion,
      platform: newPlatform,
      live: newLive,
      inactive: newInac,
    }

    members.value.data = vtuber_data
    return true
  }

  const searchMembers = (keyword) => {
    keyword = keyword.toLowerCase()

    members.value.query = keyword

    if (!keyword) {
      members.value.searchedData = []
      return
    }

    members.value.searchedData = [...toRaw(members.value.filteredData)].filter(
      (vt) => {
        const EnName = vt.EnName.toLowerCase().includes(keyword ? keyword : "")
        let JpName = vt.JpName?.includes(keyword ? keyword : "")

        return EnName || JpName || false
      }
    )

    console.log(
      `Search result from keyword "${keyword}": ${members.value.searchedData.length}`
    )
  }

  const filterMembers = () => {
    const { reg, plat, liveplat, nolive, inac } = parse(
      window.location.href,
      true
    ).query

    const regions = reg?.toLowerCase().split(",") || null

    filter.value.region = regions
    filter.value.platform = plat ? plat : null
    filter.value.live = liveplat ? liveplat : null
    filter.value.live_only = nolive?.toLowerCase() == "false" ? true : false
    filter.value.inactive = inac ? inac.toLowerCase() === "true" : false

    let { region, platform, live, live_only, inactive } = filter.value
    const regMenu = menuFilter.value.region.map((r) => r.toLowerCase())

    let vtuber_data = [...toRaw(members.value.data)]

    if (
      region ||
      platform ||
      live ||
      live_only !== false ||
      inactive !== false
    ) {
      // filter by region
      if (region) {
        vtuber_data = vtuber_data.filter((vtuber) => {
          if (regions.find((r) => !regMenu.includes(r))) return false
          return region.includes(vtuber.Region.toLowerCase())
        })
      }

      // filter by platform
      if (platform) {
        switch (platform) {
          case "yt":
            vtuber_data = vtuber_data.filter((vt) => vt.Youtube)
            break
          case "tw":
            vtuber_data = vtuber_data.filter((vt) => vt.Twitch)
            break
          case "bl":
            vtuber_data = vtuber_data.filter((vt) => vt.BiliBili)
            break
          default:
            vtuber_data = vtuber_data.filter((vt) => {
              if (platform === "-yt")
                return vt.Youtube && !vt.Twitch && !vt.BiliBili
              else if (platform === "-tw")
                return vt.Twitch && !vt.Youtube && !vt.BiliBili
              else if (platform === "-bl")
                return vt.BiliBili && !vt.Twitch && !vt.Youtube
              else if (platform.match(/^-(yt,tw|tw,yt)/g)) return !vt.BiliBili
              else if (platform.match(/^-(tw,bl|bl,tw)/g)) return !vt.Youtube
              else if (platform.match(/^-(yt,bl|bl,yt)/g)) return !vt.Twitch
              else return false
            })
            break
        }
      }

      // filter by live
      vtuber_data = vtuber_data.filter(({ IsLive }) => {
        const { Youtube, Twitch, BiliBili } = IsLive
        let live_filter = Youtube || Twitch || BiliBili

        if (live) {
          switch (live) {
            case "yt":
              live_filter = Youtube
              break
            case "tw":
              live_filter = Twitch
              break
            case "bl":
              live_filter = BiliBili
              break
          }

          if (live.match(/(yt,tw|tw,yt)/g)) live_filter = !BiliBili
          else if (live.match(/(tw,bl|bl,tw)/g)) live_filter = !Youtube
          else if (live.match(/(yt,bl|bl,yt)/g)) live_filter = !Twitch
        }

        const noLive = !BiliBili && !Youtube && !Twitch

        return live_only !== true ? live_filter || noLive : live_filter
      })

      vtuber_data = vtuber_data.filter(({ Status }) => {
        if (inactive) return Status === "Inactive"
        else return true
      })
    }

    console.log(`Total member after filtering: ${vtuber_data.length}`)

    let newPlatform = []
    let newTwitter = false

    for (const vt of vtuber_data) {
      if (!newPlatform.includes("youtube") && vt.Youtube)
        newPlatform.push("youtube")
      if (!newPlatform.includes("twitch") && vt.Twitch)
        newPlatform.push("twitch")
      if (!newPlatform.includes("bilibili") && vt.BiliBili)
        newPlatform.push("bilibili")

      if (vt.Twitter) newTwitter = true
    }

    sortMenu.value = {
      platform: newPlatform,
      twitter: newTwitter,
    }

    members.value.filteredData = vtuber_data
  }

  const sortingMembers = () => {
    members.value.query = ""
    members.value.searchedData = []

    const {
      sort,
      live: liveLink,
      inaclast,
    } = parse(window.location.href, true).query

    // if link is not /vtubers
    if (!window.location.pathname.match(/\/vtubers/g)) return

    sorting.value.type = sort ? sort?.replace("-", "") : "name"
    sorting.value.order = sort?.includes("-") ? "desc" : "asc"
    sorting.value.live = liveLink === "false" ? false : true
    sorting.value.inactive = inaclast === "false" ? false : true

    let { type, order, live, inactive } = sorting.value

    const vtuber_data = [...toRaw(members.value.filteredData)]

    vtuber_data.sort(
      // sort by name
      ({ EnName: nameA }, { EnName: nameB }) =>
        nameA.toLowerCase() < nameB.toLowerCase() ? -1 : 1
    )

    if (sort) {
      if (type === "yt") {
        console.log("Sort by youtube subscribers")
        vtuber_data.sort(({ Youtube: ytA }, { Youtube: ytB }) => {
          const subsA = ytA ? ytA.Subscriber : 0
          const subsB = ytB ? ytB.Subscriber : 0

          return subsB - subsA
        })
      }

      if (type === "bl") {
        console.log("Sort by bilibili followers")
        vtuber_data.sort(({ BiliBili: blA }, { BiliBili: blB }) => {
          const subsA = blA ? blA.Followers : 0
          const subsB = blB ? blB.Followers : 0

          return subsB - subsA
        })
      }

      if (type === "tw") {
        console.log("Sort by Twitch followers")
        vtuber_data.sort(({ Twitch: twA }, { Twitch: twB }) => {
          const subsA = twA ? twA.Followers : 0
          const subsB = twB ? twB.Followers : 0

          return subsB - subsA
        })
      }

      if (type === "twr") {
        console.log("Sort by Twitter followers")
        vtuber_data.sort(({ Twitter: twA }, { Twitter: twB }) => {
          const subsA = twA ? twA.Followers : 0
          const subsB = twB ? twB.Followers : 0

          return subsB - subsA
        })
      }

      if (type === "ytv") {
        console.log("Sort by Youtube views")
        vtuber_data.sort(({ Youtube: ytA }, { Youtube: ytB }) => {
          const viewsA = ytA ? ytA.ViwersCount : 0
          const viewsB = ytB ? ytB.ViwersCount : 0

          return viewsB - viewsA
        })
      }

      if (type === "blv") {
        console.log("Sort by bilibili views")
        vtuber_data.sort(({ BiliBili: blA }, { BiliBili: blB }) => {
          const viewsA = blA ? blA.ViwersCount : 0
          const viewsB = blB ? blB.ViwersCount : 0

          return viewsB - viewsA
        })
      }

      if (order === "desc") {
        console.log("Sort descending")
        vtuber_data.reverse()
      }
    }

    if (live) {
      console.log("Sort by live status")

      vtuber_data.sort(({ IsLive: liveA }, { IsLive: liveB }) => {
        if (!liveA.Youtube && liveB.Youtube) return 1
        if (liveA.Youtube && !liveB.Youtube) return -1
        if (!liveA.Twitch && liveB.Twitch) return 1
        if (liveA.Twitch && !liveB.Twitch) return -1
        if (!liveA.BiliBili && liveB.BiliBili) return 1
        if (liveA.BiliBili && !liveB.BiliBili) return -1
        return 0
      })
    }

    if (inactive) {
      vtuber_data.sort(({ Status: statA }, { Status: statB }) => {
        if (statA === "Inactive" && statB !== "Inactive") return 1
        if (statA !== "Inactive" && statB === "Inactive") return -1
        return 0
      })
    }

    members.value.filteredData = vtuber_data
  }

  return {
    members,
    menuFilter,
    sortMenu,
    filter,
    fetchMembers,
    filterMembers,
    sortingMembers,
    searchMembers,
  }
})
