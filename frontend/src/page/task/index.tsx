import { listTasks } from "@/api/processes"
import Loading from "@/components/common/loading"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from "@/components/ui/table"
import { cn } from "@/lib/utils"
import type { Task } from "@/model/processes"
import type { TaskStatus } from "@/model/projects"
import { getUserID } from "@/util/auth"
import { useQuery } from "@tanstack/react-query"
import { useNavigate } from "react-router"

function TaskOverview() {
  return (
    <Card className="hover:shadow-lg transition-shadow duration-300">
      <CardHeader>
        <CardTitle>概览</CardTitle>
      </CardHeader>
      <CardContent>内容</CardContent>
    </Card>
  )
}

function Task({
  taskStatus,
  className
}: {
  taskStatus: TaskStatus
  className?: string
}) {
  const userID = getUserID()
  const { data, isError, isLoading } = useQuery<Task[], Error>({
    queryKey: ["tasks", taskStatus, { userID }],
    queryFn: () => {
      if (!userID) {
        throw new Error("缺少userID")
      }
      const res = listTasks(+userID, taskStatus)
      return res
    }
  })

  const navigate = useNavigate()

  if (isLoading) {
    return (
      <div className={cn("flex items-center justify-center h-full", className)}>
        <Loading content="任务加载中..." />
      </div>
    )
  }

  if (isError) {
    return (
      <div className={cn("flex items-center justify-center h-full", className)}>
        加载失败
      </div>
    )
  }

  let title = ""
  switch (taskStatus) {
    case "todo":
      title = "待办"
      break
    case "ready":
      title = "准备"
      break
    case "doing":
      title = "进行"
      break
  }

  return (
    <Card
      className={cn(
        "hover:shadow-lg transition-shadow duration-300",
        className
      )}
    >
      <CardHeader>
        <CardTitle>{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>任务</TableHead>
              <TableHead>详情</TableHead>
              <TableHead>项目</TableHead>
              {taskStatus === "doing" && <TableHead>剩余天数</TableHead>}
            </TableRow>
          </TableHeader>
          <TableBody>
            {data?.map((t) => (
              <TableRow
                key={t.id}
                onClick={() => navigate(`/projects/${t.projectID}`)}
                className="cursor-pointer"
              >
                <TableCell>{t.name}</TableCell>
                <TableCell>{t.description}</TableCell>
                <TableCell>{t.projectName}</TableCell>
                {taskStatus === "doing" && (
                  <TableCell>{t.remainingDays}</TableCell>
                )}
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  )
}

export default function TaskPage() {
  return (
    <main className="flex flex-col gap-4 p-4 h-full">
      <div className="w-full">
        <TaskOverview />
      </div>
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <Task taskStatus="doing" />
        <Task taskStatus="ready" className="hidden sm:flex" />
        <Task taskStatus="todo" className="hidden sm:flex" />
      </div>
    </main>
  )
}
