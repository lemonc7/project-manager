import type { ListProcessesRequest, Process } from "@/model/processes"
import service from "./request"

export async function listProcesses(req: ListProcessesRequest) {
  const res = await service.post<Process[]>("/process/filter", req)
  return Array.isArray(res.data) ? res.data : []
}

export async function getProcessesByProjectID(id: number) {
  return listProcesses({ projectID: id })
}

export async function createProcess(req: Process) {
  return service.post("/process", req)
}

export async function updateProcess(id: number, req: Process) {
  return service.put(`/process/${id}`, req)
}

export async function deleteProcess(id: number) {
  return service.delete(`/process/${id}`)
}
