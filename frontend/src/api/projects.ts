import type {
  CreateProjectRequest,
  ListProjectsRequest,
  Project,
  UpdateProjectRequest
} from "@/model/projects"
import service from "./request"
import type { Process } from "@/model/processes"

export async function listProjectsFilter(req: ListProjectsRequest) {
  const res = await service.post<Project[]>("/project/filter", req)
  return Array.isArray(res.data) ? res.data : []
}

export async function listProjects() {
  const res = await service.post<Project[]>("/project/filter")
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

export async function updateProject(
  id: number,
  req: Omit<UpdateProjectRequest, "id">
) {
  const res = await service.put<Project>(`/project/${id}`, req)
  return res.data
}

export async function deleteProject(id: number) {
  return service.delete(`/project/${id}`)
}
