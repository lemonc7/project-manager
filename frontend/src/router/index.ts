import Layout from "@/components/layout/Layout.vue"
import NotFound from "@/page/404/NotFound.vue"
import LoginPage from "@/page/login/Login.vue"
import TaskPage from "@/page/task/Task.vue"
import { isLoggedIn } from "@/util/auth"
import { createRouter, createWebHistory, type RouteRecordRaw } from "vue-router"

const routes: RouteRecordRaw[] = [
  {
    path: "/login",
    component: LoginPage
  },
  {
    path: "/",
    component: Layout,
    meta: { requiresAuth: true },
    children: [
      {
        path: "",
        redirect: "/home"
      },
      {
        path: "home",
        component: TaskPage,
        meta: { title: "任务总览" }
      },
      // {
      //   path: "projects",
      //   component: ProjectLayout,
      //   children: [
      //     {
      //       path: "",
      //       component: ProjectPage,
      //       meta: { title: "项目总览" }
      //     },
      //     {
      //       path: ":id",
      //       component: ProcessPage,
      //       meta: { title: "项目流程" }
      //     }
      //   ]
      // },
      // {
      //   path: "employees",
      //   component: EmployeePage,
      //   meta: { title: "成员列表" }
      // },
      {
        path: "/:pathMatch(.*)*",
        component: NotFound,
        meta: { title: "无此页面" }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 导航守卫实现 ProtectElement 的逻辑
router.beforeEach((to, _, next) => {
  if (to.meta.requiresAuth) {
    if (!isLoggedIn()) {
      next("/login")
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router
