import Vue from 'vue'
import Router from 'vue-router'
import Guilds from '../components/Guilds.vue'
import Channels from '../components/Channels.vue'
import Channel from '../components/Channel.vue'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      beforeEnter(to, from, next) {
        window.location.href = "http://localhost:5000";
      }
    },
    {
      path: '/guilds',
      name: 'Guilds',
      component: Guilds
    },
    {
      path: '/guilds/:id/channels',
      name: 'Channels',
      component: Channels
    },
    {
      path: '/channel/:id',
      name: 'Channel',
      component: Channel
    }        
  ]
})
