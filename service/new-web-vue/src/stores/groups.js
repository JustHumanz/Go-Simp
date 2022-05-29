import { defineStore } from "pinia"
import { ref, computed } from "vue"
import axios from "axios"
import Config from "../config.json"

export const useGroupStore = defineStore("groups", () => {
  const groups = ref({
    error: false,
    status: "",
    data: [],
  })

  const id = ref({ id: null })

  const singleGroup = computed(() =>
    groups.value.data.find((group) => group.id === id.value.id)
  )

  const getId = (idGroup) => {
    id.value.id = idGroup
  }

  // const getGroups = () => computed(() => groups)

  const fetchGroups = async () => {
    let err = false
    groups.value.status = ""
    groups.value.error = false

    const data_groups = await axios
      .get(Config.REST_API + "/v2/groups/")
      .then((response) => response.data)
      .catch((error) => {
        err = true
        groups.value.error = true
        groups.value.status = error.response.status
      })

    if (err) return false

    // sort group data from GroupName
    data_groups.sort((a, b) => {
      if (a.GroupName.toLowerCase() < b.GroupName.toLowerCase()) return -1
      if (a.GroupName.toLowerCase() > b.GroupName.toLowerCase()) return 1
      return 0
    })

    groups.value.data = data_groups
    console.log(`Total group: ${data_groups.length}`)
  }

  return {
    groups,
    getId,
    singleGroup,
    fetchGroups,
  }
})
