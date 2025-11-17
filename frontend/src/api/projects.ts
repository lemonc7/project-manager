import type {
  CreateProjectRequest,
  ListProjectsRequest,
  Project,
  UpdateProjectRequest
} from "@/model/projects"
import service from "./request"

export async function listProjectsFilter(req: ListProjectsRequest) {
  const res = await service.post<Project[]>("/project/filter", req)
  return Array.isArray(res.data) ? res.data : []
}

export async function listProjects() {
  const res = await service.get<Project[]>("/project/filter")
  return Array.isArray(res.data) ? res.data : []
}

export async function createProject(req: CreateProjectRequest) {
  const res = await service.post<Project[]>("/project/filter", req)
  return Array.isArray(res.data) ? res.data : []
}

export async function updateProject(
  id: number,
  req: Omit<UpdateProjectRequest, "id">
) {
  return service.put(`/project/${id}`, req)
}

export async function deleteProject(id: number) {
  return service.delete(`/project/${id}`)
}
