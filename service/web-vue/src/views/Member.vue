<template>
  <div id="app">
    <div class="container text-center" id="container" v-for:="(Member, index) in Members">
      <div v-if="Member.Youtube && Member.BiliBili != null" id="carouselExampleSlidesOnly" class="carousel slide" data-ride="carousel">
        <div class="carousel-inner">
          <div class="carousel-item active">
            <a v-bind:href="'https://www.youtube.com/channel/'+Member.Youtube.ID+'?sub_confirmation=1'" target="_blank" rel="noopener noreferrer">
              <img class="w-20 rounded mx-auto" width="200" height="200" v-bind:src=Member.Youtube.Avatar @error="replaceByDefault" alt="Youtube">
            </a>
          </div>
          <div class="carousel-item">
            <a v-bind:href='"http://space.bilibili.com/" + Member.BiliBili.ID' target="_blank" rel="noopener noreferrer">
              <img class="w-20 rounded mx-auto" width="200" height="200" v-bind:src=Member.BiliBili.Avatar @error="replaceByDefault" alt="BiliBili">
            </a>
          </div>
        </div>
      </div>
      <div v-else-if="Member.Youtube != null && Member.BiliBili == null">
        <a v-bind:href="'https://www.youtube.com/channel/'+ Member.Youtube.ID+'?sub_confirmation=1'" target="_blank" rel="noopener noreferrer">
          <img class="w-20 rounded mx-auto" width="200" height="200" :src=Member.Youtube.Avatar @error="replaceByDefault" alt="Youtube"></a>
      </div>
      <div v-else-if="Member.Youtube == null && Member.BiliBili != null">
        <a v-bind:href='"http://space.bilibili.com/" + Member.BiliBili.ID' target="_blank" rel="noopener noreferrer">
          <img class="w-20 rounded mx-auto" width="200" height="200" :src=Member.BiliBili.Avatar @error="replaceByDefault" alt="BiliBili"></a>
      </div>
      <div class="container mb-5">
        <div class="row">
           <div class="col text-center pt-2">
             <p v-if="Member.EnName && Member.JpName != null" >Original Name: {{Member.EnName}}/{{Member.JpName}}</p>
             <p v-else-if="Member.EnName != null && Member.JpName == null">Original Name: {{Member.EnName}}</p>
             <p v-else-if="Member.EnName == null && Member.JpName != null" >Original Name: {{Member.JpName}}</p>
             <p>Nickname: {{Member.NickName}}</p>
             <p>Region: {{Member.Region}}</p>
           </div>
        </div>
        <div class="row text-center d-flex justify-content-center">
           <div v-if="Member.Youtube != null" class="col-xs-6 col-sm-4 col-md-4 mt-4 col">
              <a v-bind:href="'https://www.youtube.com/channel/'+Member.Youtube.ID+'?sub_confirmation=1'" target="_blank" rel="noopener noreferrer">
                <font-awesome-icon :icon="{ prefix: 'fab', iconName: 'youtube' }" size="2x" style="color: #FF0000"/></a>
                <h2><CountTo :endVal=Member.Youtube.Subscriber></CountTo></h2>
                <p class="count-text ">Youtube subscriber</p>
           </div>
           <div v-if="Member.Twitter != null" class="col-xs-6 col-sm-4 col-md-4 mt-4 col">
              <a v-bind:href="'https://twitter.com/' + Member.Twitter.UserName" target="_blank" rel="noopener noreferrer">
                <font-awesome-icon :icon="{ prefix: 'fab', iconName: 'twitter' }" size="2x" style="color: #1DA1F2"/></a>
                <h2><CountTo :endVal=Member.Twitter.Followers></CountTo></h2>
                <p class="count-text ">Twitter followers</p>
           </div>
           <div v-if="Member.BiliBili != null" class="col-xs-6 col-sm-4 col-md-4 mt-4 col">
              <a v-bind:href="'https://space.bilibili.com/' + Member.BiliBili.ID" target="_blank" rel="noopener noreferrer">
              <svg class="icon" style="width: 2em; height: 2em;vertical-align: middle;fill: #23aee5;overflow: hidden;" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="517"><path d="M777.514667 131.669333a53.333333 53.333333 0 0 1 0 75.434667L728.746667 255.829333h49.92A160 160 0 0 1 938.666667 415.872v320a160 160 0 0 1-160 160H245.333333A160 160 0 0 1 85.333333 735.872v-320a160 160 0 0 1 160-160h49.749334L246.4 207.146667a53.333333 53.333333 0 1 1 75.392-75.434667l113.152 113.152c3.370667 3.370667 6.186667 7.04 8.448 10.965333h137.088c2.261333-3.925333 5.12-7.68 8.490667-11.008l113.109333-113.152a53.333333 53.333333 0 0 1 75.434667 0z m1.152 231.253334H245.333333a53.333333 53.333333 0 0 0-53.205333 49.365333l-0.128 4.010667v320c0 28.117333 21.76 51.157333 49.365333 53.162666l3.968 0.170667h533.333334a53.333333 53.333333 0 0 0 53.205333-49.365333l0.128-3.968v-320c0-29.44-23.893333-53.333333-53.333333-53.333334z m-426.666667 106.666666c29.44 0 53.333333 23.893333 53.333333 53.333334v53.333333a53.333333 53.333333 0 1 1-106.666666 0v-53.333333c0-29.44 23.893333-53.333333 53.333333-53.333334z m320 0c29.44 0 53.333333 23.893333 53.333333 53.333334v53.333333a53.333333 53.333333 0 1 1-106.666666 0v-53.333333c0-29.44 23.893333-53.333333 53.333333-53.333334z" p-id="518"></path></svg></a>
                <h2><CountTo :endVal=Member.BiliBili.Followers></CountTo></h2>
                <p class="count-text ">BiliBili followers</p>
           </div>
           <div v-if="Member.Twitch != null" class="col-xs-6 col-sm-4 col-md-4 mt-4 col">
              <a v-bind:href="'https://www.twitch.tv/' + Member.Twitch.UserName" target="_blank" rel="noopener noreferrer">
              <font-awesome-icon :icon="{ prefix: 'fab', iconName: 'twitch' }" size="2x" style="color: #6441a5"/></a>
                <h2><CountTo :endVal=Member.Twitch.Followers></CountTo></h2>
                <p class="count-text ">Twitch followers</p>
           </div>           
        </div>
      </div> 
    </div>
  </div>
</template>

<script>
import { CountTo } from 'vue3-count-to';
import axios from 'axios'
import Config from "../config.json";
export default {
  name: 'Vtuber',
  data () {
    return {
      Members: null
    }
  },
  components: {
    CountTo
  },
  mounted () {
    axios
      .get(Config.REST_API+'/members/' + this.$route.params.id).then(response => (this.Members = response.data))
  },
  methods: {
    replaceByDefault(e) {
        e.target.src = require('../assets/404.jpg')
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
