import type { InsertProcessRequest, ListProcessesRequest, Process } from "@/model/processes"
import service from "./request"

export async function listProcesses(req: ListProcessesRequest) {
  const res = await service.post<Process[]>("/process/filter", req)
  return Array.isArray(res.data) ? res.data : []
}

export async function getProcessesByProjectID(id: number) {
  return listProcesses({ projectID: id })
}

export async function insertProcessLast(req: Omit<InsertProcessRequest,"sequence">) {
  return service.post("/process", req)
}

export async function insertProcessAt(req: InsertProcessRequest) {
  const {sequence, ...rest} = req
  return service.post(`/process/${sequence}`, rest)
}

export async function updateProcess(id: number, req: Process) {
  return service.put(`/process/${id}`, req)
}

export async function deleteProcess(id: number) {
  return service.delete(`/process/${id}`)
}
