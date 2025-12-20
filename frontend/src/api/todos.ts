import type {
  CreateTodoRequest,
  Todo,
  TodoStatus,
  UpdateTodoRequest
} from "../models/todos"
import service from "./request"

export async function listTodosByReporter(
  reporterID: number,
  status?: TodoStatus[]
) {
  const res = await service.get<Todo[]>(`/todo/reporter/${reporterID}`, {
    params: {
      status: status
    }
  })
  return res.data
}

export async function listTodosByAssignee(
  assigneeID: number,
  status?: TodoStatus[]
) {
  const res = await service.get<Todo[]>(`/todo/assignee/${assigneeID}`, {
    params: {
      status: status
    }
  })
  return res.data
}

export async function createTodo(req: CreateTodoRequest) {
  return service.post("/todo", req)
}

export async function updateTodo(id: number, req: UpdateTodoRequest) {
  return service.patch(`/todo/${id}`, req)
}

export async function deleteTodo(id: number) {
  return service.patch(`/todo/${id}`)
}
