<template>
  <div class="search-wrapper p-3 d-flex flex-row-reverse">
    <input type="text" v-model="search" :placeholder="ExName + '/'+ ExGroup" />
  </div>
    <div class="container text-center wrapper" id="container">
        <router-link v-for='m in filteredList' :key='m.ID' :to="`/Member/${m.ID}`">
        <figure class="figure p-2 text-center card-search" :name="m.Region" :id="m.NickName">
            <img v-if="m.Youtube != null" :src=m.Youtube.Avatar class="zoom figure-img img-fluid img-thumbnail" alt="Youtube" width="75" height="75"  @error="replaceByDefault">
            <img v-else-if="m.BiliBili != null" :src=m.BiliBili.Avatar class="zoom figure-img img-fluid img-thumbnail" alt="BiliBili" width="75" height="75"  @error="replaceByDefault">
            <img v-else-if="m.Twitch != null" :src=m.Twitch.Avatar class="zoom figure-img img-fluid img-thumbnail" alt="Twitch" width="75" height="75"  @error="replaceByDefault">
            <img v-else src=https://cdn.humanz.moe/404.jpg class="zoom figure-img img-fluid img-thumbnail" alt="Youtube" width="75" height="75">
         <figcaption class="figure-caption">{{m.EnName}}</figcaption>
         </figure>
        </router-link>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  name: 'Members',
  data () {
    return {
      Members: [],
      search: '',
      ExName: '',
      ExGroup: ''
    }
  },

  methods: {
    replaceByDefault(e) {
        e.target.src = require('../assets/404.jpg')
    },
    getData(){
      axios.get(process.env.VUE_APP_RESTAPI+'/members/').then(response => {
      for (let i = 0; i < response.data.length; i++) {
        if (response.data[i]["Youtube"] != null) {
          response.data[i]["Youtube"]["Avatar"]= response.data[i]["Youtube"]["Avatar"].replace("s800","s75")
        }
      }
      this.Members = response.data
      var tmp = this.Members[Math.floor(Math.random() * this.Members.length)]
      this.ExName = tmp["EnName"]
      this.ExGroup = tmp["GroupName"]
    }) 
    }
  },
  computed: {
    filteredList: function() {
      return this.Members.filter(post => {
        let EnName = post.EnName.toLowerCase().includes(this.search.toLowerCase())
        let JpName 
        if (post.JpName != null) {
          JpName  = post.JpName.toLowerCase().includes(this.search.toLowerCase())
        } 
        let NickName = post.NickName.toLowerCase().includes(this.search.toLowerCase())

        if (this.search.length == 2) {
          let Region = post.Region.toLowerCase().includes(this.search.toLowerCase())
          return Region
        }

        let Group = post.GroupName.toLowerCase().includes(this.search.toLowerCase())
        return EnName || JpName || NickName || Group
      })
    }
  },
  mounted () {
    this.getData()
  },        
}

</script>
