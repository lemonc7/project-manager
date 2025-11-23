import NotFoundPage from "@/pages/404"
import TaskPage from "@/pages/task"
import CustomLayout from "@/components/layout/custom-layout"
import { createBrowserRouter, redirect } from "react-router"
import ProjectPage from "@/pages/project"
import ProjectLayout from "@/components/layout/project-layout"
import ProcessPage from "@/pages/process"
import EmployeePage from "@/pages/employee"

const router = createBrowserRouter([
  {
    path: "/",
    element: <CustomLayout />,
    children: [
      {
        index: true,
        loader: () => redirect("/home")
      },
      {
        path: "home",
        element: <TaskPage />,
        handle: { title: "任务总览" }
      },
      {
        path: "projects",
        element: <ProjectLayout />,
        children: [
          {
            index: true,
            element: <ProjectPage />,
            handle: { title: "项目总览" }
          },
          {
            path: ":id",
            element: <ProcessPage />,
            handle: { title: "项目流程" }
          }
        ]
      },
      {
        path: "employees",
        element: <EmployeePage />,
        handle: { title: "成员列表" }
      },
      {
        path: "ping",
        element: <h1>pong</h1>
      },
      {
        path: "*",
        element: <NotFoundPage />,
        handle: { title: "无此页面" }
      }
    ]
  }
])

export default router
