import Vue from 'vue'
import Router from 'vue-router'
import Group from '@/components/Group'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Index'
    },
    {
      path: '/group/:id',
      name: 'Group',
      component: Group
    }
  ]
})
