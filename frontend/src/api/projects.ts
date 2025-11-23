import type {
  CreateProjectRequest,
  ListProjectsRequest,
  Project,
  UpdateProjectRequest
} from "@/models/projects"
import service from "./request"
import type { Process } from "@/models/processes"

export async function listProjectsFilter(req: ListProjectsRequest) {
  const res = await service.get<Project[]>("/project/filter", { params: req })
  return Array.isArray(res.data) ? res.data : []
}

export async function listProcesses(id: number) {
  const res = await service.get<Process[]>(`/project/${id}`)
  return Array.isArray(res.data) ? res.data : []
}

export async function createProject(req: CreateProjectRequest) {
  const res = await service.post<Project>("/project", req)
  return res.data
}

export async function updateProject(id: number, req: UpdateProjectRequest) {
  const res = await service.put<Project>(`/project/${id}`, req)
  return res.data
}

export async function deleteProject(id: number) {
  return service.delete(`/project/${id}`)
}
