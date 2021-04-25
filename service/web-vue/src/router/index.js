import { createRouter, createWebHashHistory } from 'vue-router'
import Index from '../views/Index.vue'
import Group from '../views/Group.vue'
import Members from '../views/Members.vue'
import Vtuber from '../views/Member.vue'
import Exec from '../views/Exec.vue'
import Guide from '../views/Guide.vue'


const routes = [
  {
    path: '/',
    name: 'Index',
    component: Index
  },
  {
    path: '/Group/:id',
    name: 'Group',
    component: Group
  },
  {
    path: '/Vtubers/',
    name: 'Members',
    component: Members
  },
  {
    path: '/Member/:id',
    name: 'Vtuber',
    component: Vtuber
  },
  {
    path: '/Exec/',
    name: 'Exec',
    component: Exec
  },
  {
    path: '/Guide/',
    name: 'Guide',
    component: Guide
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
