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

// 导航守卫实现 requiresAuth 的逻辑
router.beforeEach(async (to) => {
  const loggedIn = isLoggedIn()
  const requiresAuth = to.matched.some((record) => record.meta.requiresAuth)
  // 已登录，重定向到task页面
  if (to.name === "login" && loggedIn) {
    return { name: "task", replace: true }
  }
  // 未登录，重定向到login页面
  if (requiresAuth && !loggedIn) {
    return { name: "login" }
  }
  // 其他情况正常放行
  return true
})

let projectScrollTop = 0
// 离开 /project 时保存滚动位置
router.beforeEach((_, from) => {
  if (from.path === "/project") {
    const table = document.getElementById("project-table")
    if (table) {
      projectScrollTop = table.scrollTop
    }
  }
})

// 进入 /project 时恢复滚动位置
router.afterEach((to) => {
  if (to.path === "/project") {
    setTimeout(() => {
      const table = document.getElementById("project-table")
      if (table) {
        table.scrollTop = projectScrollTop
      }
    }, 0)
  }
})


export default router
