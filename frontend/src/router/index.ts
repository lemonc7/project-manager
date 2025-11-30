import Layout from "@/components/layout/Layout.vue"
import ProjectLayout from "@/components/layout/ProjectLayout.vue"
import NotFoundPage from "@/pages/404/NotFoundPage.vue"
import EmployeePage from "@/pages/employee/EmployeePage.vue"
import LoginPage from "@/pages/login/LoginPage.vue"
import ProcessPage from "@/pages/process/ProcessPage.vue"
import ProjectPage from "@/pages/project/ProjectPage.vue"
import TaskPage from "@/pages/task/TaskPage.vue"
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
      {
        path: "employee",
        component: EmployeePage,
        meta: { title: "成员列表" }
      },
      {
        path: "project",
        component: ProjectLayout,
        children: [
          {
            path: "",
            component: ProjectPage,
            meta: { title: "项目总览" }
          },
          {
            path: ":id",
            component: ProcessPage,
            meta: { title: "项目流程" }
          }
        ]
      },
      {
        path: "/:pathMatch(.*)*",
        component: NotFoundPage,
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
