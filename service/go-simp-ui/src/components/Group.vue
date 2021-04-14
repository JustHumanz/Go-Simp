<template>
  <div class="go-simp-group">
  <ul>
    <li v-for='m in msg' :key='m._id'>
      {{ m.NickName }}
    </li>
  </ul>
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
      .get('http://localhost:2525/members/13,14', { crossDomain: true })
      .then(response => (this.msg = response.data))
  }
}

</script>
