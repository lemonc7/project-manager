import { isLoggedIn } from "../utils/auth"
import {
  createRouter,
  createWebHistory,
  RouterView,
  type RouteRecordRaw
} from "vue-router"

const routes: RouteRecordRaw[] = [
  {
    path: "/login",
    name: "login",
    component: () => import("../pages/LoginPage.vue")
  },
  {
    path: "/",
    component: () => import("../components/layout/index.vue"),
    meta: { requiresAuth: true },
    children: [
      {
        path: "",
        name: "task",
        component: () => import("../pages/TaskPage.vue"),
        meta: { title: "任务总览" }
      },
      {
        path: "employee",
        name: "employee",
        component: () => import("../pages/EmployeePage.vue"),
        meta: { title: "成员列表" }
      },
      {
        path: "team",
        name: "team",
        component: () => import("../pages/TeamPage.vue"),
        meta: { title: "小组列表" }
      },
      {
        path: "project",
        name: "project",
        component: () => import("../pages/ProjectPage.vue"),
        meta: { title: "项目总览" }
      },
      {
        path: "project/:id",
        name: "process",
        component: () => import("../pages/ProcessPage.vue"),
        meta: { title: "项目流程" }
      },
      {
        path: "/:pathMatch(.*)*",
        component: () => import("../pages/NotFoundPage.vue"),
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
  const requiresAuth = to.matched.some((record) => record.meta.requiresAuth)
  // 已登录，重定向到/页面
  if (to.path === "/login" && isLoggedIn()) {
    return next("/")
  }
  // 未登录，重定向到login页面
  if (requiresAuth && !isLoggedIn()) {
    return next("/login")
  }
  next()
})

export default router
