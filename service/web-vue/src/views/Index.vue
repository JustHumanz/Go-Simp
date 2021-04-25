<template>
  <div id="app">  
    <div class="container" id="container">
      <div class="text-center">
        <figure v-for='m in msg' :key='m._id' class="figure p-4">
          <router-link :to="`/group/${m.ID}`">
            <img v-bind:src=m.GroupIcon class="figure-img img-fluid" alt="" width="200" height="200"/>
          </router-link>
            <figcaption class="figure-caption" style="text-transform: capitalize;">{{ m.GroupName }}</figcaption>
        </figure>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  name: 'Index',
  data () {
    return {
      msg: null
    }
  },
  mounted () {   
    axios.get(process.env.VUE_APP_RESTAPI+'/groups/').then(response => (this.msg = response.data))
  }
}
</script>

<style>
#app {
  min-height: 75rem;
  padding-top: 4.5rem;
}
</style>
