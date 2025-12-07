import { isLoggedIn } from "../utils/auth"
import { createRouter, createWebHistory, type RouteRecordRaw } from "vue-router"

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
        component: () => import("../pages/TaskPage.vue")
      },
      {
        path: "employee",
        name: "employee",
        component: () => import("../pages/EmployeePage.vue")
      },
      {
        path: "team",
        name: "team",
        component: () => import("../pages/TeamPage.vue")
      },
      {
        path: "project",
        name: "project",
        component: () => import("../pages/ProjectPage.vue")
      },
      {
        path: "project/:id",
        name: "process",
        component: () => import("../pages/ProcessPage.vue")
      },
      {
        path: "/:pathMatch(.*)*",
        component: () => import("../pages/NotFoundPage.vue")
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
