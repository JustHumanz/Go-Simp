<template>
  <div id="app">
    <div class="container" id="container">
      <div class="text-center">
        <div class="row row-cols-1 row-cols-md-4">
            <div class="col mb-4" v-for='m in msg' :key='m._id'>
                <div class="card" v-bind:id=m.NickName>
                    <div v-if="m.Youtube !== null">
                        <img class="card-img-top" v-bind:src=m.Youtube.Avatar alt="Card image cap">
                    </div>
                    <div v-else-if="m.BiliBili !== 'B'">
                        <img class="card-img-top" v-bind:src=m.BiliBili.Avatar alt="Card image cap">
                    </div>
                    <div v-else>
                        <img class="card-img-top" src="https://cdn.humanz.moe/404.jpg" alt="Card image cap">
                    </div>
                    <div class="card-body">
                    <h5 class="card-title">{{Member.NickName}}</h5>
                        <router-link :to="`/Member/${m.ID}`">Info </router-link>
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
      msg: null
    }
  },
  mounted () {
    axios
      .get('http://localhost:2525/members/?groupid=' + this.$route.params.id, { crossDomain: true })
      .then(response => (this.msg = response.data))
  }
}
</script>

<style>
#app {
  min-height: 75rem;
  padding-top: 4.5rem;
}
</style>
