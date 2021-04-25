<template>
  <div class="search-wrapper p-3 d-flex flex-row-reverse">
    <input type="text" v-model="search" :placeholder=ExName />
  </div>
  <div id="app">
    <div class="container" id="container">
      <div class="text-center">
        <div class="row row-cols-1 row-cols-md-4">
            <div class="col mb-4" v-for='m in filteredList' :key='m.ID'>
                <div class="card" v-bind:id=m.NickName>
                    <div v-if="m.Youtube !== null">
                      <router-link :to="`/Member/${m.ID}`">
                        <img class="card-img-top" v-bind:src=m.Youtube.Avatar onerror="this.src='https://cdn.humanz.moe/404.jpg'" alt="Card image cap">
                      </router-link>
                    </div>
                    <div v-else-if="m.BiliBili !== null">
                      <router-link :to="`/Member/${m.ID}`">
                        <img class="card-img-top" v-bind:src=m.BiliBili.Avatar onerror="this.src='https://cdn.humanz.moe/404.jpg'" alt="Card image cap">
                      </router-link>
                    </div>
                    <div v-else>
                      <router-link :to="`/Member/${m.ID}`">
                        <img class="card-img-top" src="https://cdn.humanz.moe/404.jpg" alt="Card image cap">
                      </router-link>
                    </div>
                    <div class="card-body">
                    <h5 class="card-title">{{m.NickName}}</h5>
                      <div class="row text-center">
                      <div v-if="m.Youtube != null" class="col-xs-6 col-sm-4 col-md-4 mt-4">
                          <div class="counter-div">
                            <a :href="'https://www.youtube.com/channel/'+m.Youtube.ID+'?sub_confirmation=1'" target="_blank" rel="noopener noreferrer">
                            <font-awesome-icon :icon="{ prefix: 'fab', iconName: 'youtube' }" size="2x" style="color: #FF0000"/></a>
                          </div>
                      </div>
                      <div v-if="m.Twitter != null" class="col-xs-6 col-sm-4 col-md-4 mt-4">
                          <div class="counter-div" >
                            <a :href="'https://twitter.com/'+m.Twitter.UserName" target="_blank" rel="noopener noreferrer">
                            <font-awesome-icon :icon="{ prefix: 'fab', iconName: 'twitter' }" size="2x" style="color: #1DA1F2"/></a>
                          </div>
                      </div>
                      <div v-if="m.BiliBili != null" class="col-xs-6 col-sm-4 col-md-4 mt-4">
                          <div class="counter-div">
                            <a :href="'https://space.bilibili.com/'+m.BiliBili.ID" target="_blank" rel="noopener noreferrer">
                            <svg class="icon" style="width: 2em; height: 2em;vertical-align: middle;fill: #23aee5;overflow: hidden;" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="517"><path d="M777.514667 131.669333a53.333333 53.333333 0 0 1 0 75.434667L728.746667 255.829333h49.92A160 160 0 0 1 938.666667 415.872v320a160 160 0 0 1-160 160H245.333333A160 160 0 0 1 85.333333 735.872v-320a160 160 0 0 1 160-160h49.749334L246.4 207.146667a53.333333 53.333333 0 1 1 75.392-75.434667l113.152 113.152c3.370667 3.370667 6.186667 7.04 8.448 10.965333h137.088c2.261333-3.925333 5.12-7.68 8.490667-11.008l113.109333-113.152a53.333333 53.333333 0 0 1 75.434667 0z m1.152 231.253334H245.333333a53.333333 53.333333 0 0 0-53.205333 49.365333l-0.128 4.010667v320c0 28.117333 21.76 51.157333 49.365333 53.162666l3.968 0.170667h533.333334a53.333333 53.333333 0 0 0 53.205333-49.365333l0.128-3.968v-320c0-29.44-23.893333-53.333333-53.333333-53.333334z m-426.666667 106.666666c29.44 0 53.333333 23.893333 53.333333 53.333334v53.333333a53.333333 53.333333 0 1 1-106.666666 0v-53.333333c0-29.44 23.893333-53.333333 53.333333-53.333334z m320 0c29.44 0 53.333333 23.893333 53.333333 53.333334v53.333333a53.333333 53.333333 0 1 1-106.666666 0v-53.333333c0-29.44 23.893333-53.333333 53.333333-53.333334z" p-id="518"></path></svg></a>
                          </div>
                      </div>
                      <div v-if="m.Twitch != null" class="col-xs-6 col-sm-4 col-md-4 mt-4">
                          <div class="counter-div">
                            <a :href="'https://twitch.tv/'+m.Twitch.UserName" target="_blank" rel="noopener noreferrer">
                            <font-awesome-icon :icon="{ prefix: 'fab', iconName: 'twitch' }" size="2x" style="color: #6441a5" /></a>
                          </div>
                      </div>                      
                      </div>
                    </div>
                    <div v-if="m.IsYtLive">
                    <div class="card-footer">
                        <small class="text-muted">🔴 Live</small>
                    </div>                        
                    </div>
                </div>
            </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  name: 'Group',
  data () {
    return {
      Members: [],
      ExName: '',
      search: ''
    }
  },
  mounted () {
    axios
      .get(process.env.VUE_APP_RESTAPI+'/members/?groupid=' + this.$route.params.id, { crossDomain: true,params:{
        live: 'true'
      }})
      .then(response => {
        this.Members = response.data
        this.ExName = this.Members[Math.floor(Math.random() * this.Members.length)]["EnName"]
      })
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
        let Region = post.Region.toLowerCase().includes(this.search.toLowerCase())
        let Group = post.GroupName.toLowerCase().includes(this.search.toLowerCase())
        return EnName || JpName || NickName || Region || Group
      })
    }
  },
  methods: {
    VtuberRandom() {
        var vtran = this.Members[Math.floor(Math.random() * this.Members.length)];
        console.log(this.Members.length);
        return vtran["NickName"]
    }
  }    
}
</script>

<style>
#app {
  min-height: 75rem;
  padding-top: 4.5rem;
}
</style>
