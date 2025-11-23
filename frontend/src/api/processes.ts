import type {
  InsertProcessRequest,
  Process,
  Task,
  UpdateProcessRequest
} from "@/model/processes"
import service from "./request"
import type { TaskStatus } from "@/model/projects"

export async function listProcesses(projectID: number) {
  const res = await service.get<Process[]>(`/project/${projectID}`)
  return Array.isArray(res.data) ? res.data : []
}

export async function listTasks(id: number, status: TaskStatus) {
  const res = await service.get<Task[]>(`/process/task/${id}`, {
    params: { status: status }
  })
  return Array.isArray(res.data) ? res.data : []
}

export async function insertProcessLast(req: InsertProcessRequest) {
  const res = await service.post<Process>("/process", req)
  return res.data
}

export async function insertProcessAt(
  sequence: number,
  req: InsertProcessRequest
) {
  return service.post(`/process/${sequence}`, req)
}

export async function updateProcess(id: number, req: UpdateProcessRequest) {
  const res = await service.put<Process>(`/process/${id}`, req)
  return res.data
}

export async function deleteProcess(id: number) {
  return service.delete(`/process/${id}`)
}
