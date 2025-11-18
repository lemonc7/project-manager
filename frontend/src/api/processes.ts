import type {
  InsertProcessRequest,
  ListProcessesRequest,
  Process
} from "@/model/processes"
import service from "./request"

export async function listProcessesFilter(req: ListProcessesRequest) {
  const res = await service.post<Process[]>("/process/filter", req)
  return Array.isArray(res.data) ? res.data : []
}

export async function insertProcessLast(
  req: Omit<InsertProcessRequest, "sequence">
) {
  const res = await service.post<Process>("/process", req)
  return res.data
}

export async function insertProcessAt(req: InsertProcessRequest) {
  const { sequence, ...rest } = req
  return service.post(`/process/${sequence}`, rest)
}

export async function updateProcess(id: number, req: Process) {
  const res = await service.put<Process>(`/process/${id}`, req)
  return res.data
}

export async function deleteProcess(id: number) {
  return service.delete(`/process/${id}`)
}
