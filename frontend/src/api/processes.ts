import type {
  FinishProcessRequest,
  InsertProcessRequest,
  Process,
  ReworkProcessRequest,
  Task,
  UpdateProcessRequest
} from "../models/processes"
import service from "./request"
import type { TaskStatus } from "../models/projects"

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
  return service.post<Process>("/process", req)
}

export async function insertProcessAt(
  sequence: number,
  req: InsertProcessRequest
) {
  return service.post(`/process/${sequence}`, req)
}

export async function updateProcess(id: number, req: UpdateProcessRequest) {
  return service.patch<Process>(`/process/${id}`, req)
}

export async function deleteProcess(id: number) {
  return service.delete(`/process/${id}`)
}

export async function finishProcess(id: number, req: FinishProcessRequest) {
  return service.post(`/process/finish/${id}`, req)
}

export async function reworkProcess(id: number, req: ReworkProcessRequest) {
  return service.post(`/process/rework/${id}`, req)
}
