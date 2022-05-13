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
  },
  emits: ["back"],
  async mounted() {
    await this.addVtuber()
    console.log(this.vtubers)
  },
  methods: {
    async addVtuber(e = null) {
      if (e === null) e = document.querySelector("[name=add-vtuber]")
      if (this.vtubers.length == 14) e.target.disabled = true

      this.vtubers.push(this.vtubers.length)
    },
  },
}
