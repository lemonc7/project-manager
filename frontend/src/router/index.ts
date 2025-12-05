import Layout from "../components/layout/Layout.vue"
import ProjectLayout from "../components/layout/ProjectLayout.vue"
import EmployeePage from "../pages/EmployeePage.vue"
import LoginPage from "../pages/LoginPage.vue"
import NotFoundPage from "../pages/NotFoundPage.vue"
import ProcessPage from "../pages/ProcessPage.vue"
import ProjectPage from "../pages/ProjectPage.vue"
import TaskPage from "../pages/TaskPage.vue"
import { isLoggedIn } from "../utils/auth"
import { createRouter, createWebHistory, type RouteRecordRaw } from "vue-router"

const routes: RouteRecordRaw[] = [
  {
    path: "/login",
    name: "login",
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
        name: "taskList",
        component: TaskPage,
        meta: { title: "任务总览" }
      },
      {
        path: "employee",
        name: "employeeList",
        component: EmployeePage,
        meta: { title: "成员列表" }
      },
      {
        path: "project",
        name: "projectRoot",
        component: ProjectLayout,
        children: [
          {
            path: "",
            name: "projectList",
            component: ProjectPage,
            meta: { title: "项目总览" }
          },
          {
            path: ":id",
            name: "processList",
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
