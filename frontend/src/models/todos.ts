export const TodoStatusOptions = ["todo", "doing", "done"] as const

export type TodoStatus = (typeof TodoStatusOptions)[number]

interface CreateTodoRequest {
  name: string
  description: string
  remark: string
  priority: number
  reporterID: number
  assigneeID: number
}

interface UpdateTodoRequest {
  name?: string
  description?: string
  remark?: string
  priority?: number
  status?: TodoStatus
  reporterID?: number
  assigneeID?: number
}

interface Todo {
  id: number
  name: string
  description: string
  remark: string
  status: TodoStatus
  priority: number
  employeeID: number
  employeeName: string
  startedAt?: string
  finishedAt?: string
}

export type { CreateTodoRequest, UpdateTodoRequest, Todo }
