<template>
  <div id="app">
    <div class="container" id="container">
      <div class="text-center">
        <figure v-for='m in msg' :key='m._id' class="figure p-2">
          <router-link :to="`/group/${m.ID}`">
            <img v-bind:src=m.GroupIcon class="figure-img img-fluid img-thumbnail" alt="" width="200" height="200" />
            <figcaption class="figure-caption">{{ m.GroupName }}</figcaption>
          </router-link>
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
    axios
      .get('http://localhost:2525/groups/', { crossDomain: true })
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
