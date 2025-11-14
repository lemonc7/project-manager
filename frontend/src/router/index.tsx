import NotFoundPage from "@/pages/404"
import TaskPage from "@/pages/task"
import CustomLayout from "@/components/layout/CustomLayout"
import { createBrowserRouter, redirect } from "react-router"
import ProjectPage from "@/pages/project"
import ProjectLayout from "@/components/layout/ProjectLayout"
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
        element: <TaskPage />
      },
      {
        path: "projects",
        element: <ProjectLayout />,
        children: [
          {
            index: true,
            element: <ProjectPage />
          },
          {
            path: ":id",
            element: <ProcessPage />
          }
        ]
      },
      {
        path: "employees",
        element: <EmployeePage />
      },
      {
        path: "ping",
        element: <h1>pong</h1>
      },
      {
        path: "*",
        element: <NotFoundPage />
      }
    ]
  }
])

export default router
