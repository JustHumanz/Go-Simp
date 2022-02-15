import { createRouter, createWebHistory } from "vue-router"
import HomeView from "../views/HomeView.vue"
import VtuberView from "../views/VtuberView.vue"

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  scrollBehavior(to, from, savedPosition) {
    // always scroll to top when route changed
    return { top: 0 }
  },
  routes: [
    {
      path: "/",
      name: "Home",
      component: HomeView,
    },
    {
      path: "/vtuber",
      name: "Vtuber List",

      component: VtuberView,
    },
    {
      path: "/vtuber/:id",
      name: "Group Vtuber List",

      component: VtuberView,
    },
    //{ path: "*", component: NotFound, name: "NotFound" },
    {
      path: "/invite",
      name: "invite",
      component: () =>
        (window.location.href =
          "https://discord.com/oauth2/authorize?client_id=721964514018590802&permissions=456720&scope=bot%20applications.commands"),
    },
  ],
})

export default router
