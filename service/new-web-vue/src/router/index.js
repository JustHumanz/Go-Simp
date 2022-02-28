import { createRouter, createWebHistory } from "vue-router"
import HomeView from "../views/HomeView.vue"
import ListView from "../views/ListView.vue"
import OldListView from "../views/OldListView.vue"
import VtuberView from "../views/VtuberView.vue"
import DocsView from "../views/DocsView.vue"
import SupportView from "../views/SupportView.vue"

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
      component: ListView,
    },
    {
      path: "/vtuber/:id",
      name: "Group Vtuber List",
      component: ListView,
    },
    {
      path: "/oldvtuber",
      name: "Old Vtuber List",
      component: OldListView,
    },
    {
      path: "/oldvtuber/:id",
      name: "Old Group Vtuber List",
      component: OldListView,
    },
    {
      path: "/vtuber/members/:id",
      name: "Vtuber Details",
      component: VtuberView,
    },
    {
      path: "/docs",
      name: "Documentation",
      component: DocsView,
    },
    {
      path: "/support",
      name: "Support",
      component: SupportView,
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
