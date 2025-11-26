import type {
  CreateProjectRequest,
  ListProjectsRequest,
  Project,
  UpdateProjectRequest
} from "@/model/projects"
import service from "./request"

export async function listProjectsFilter(req: ListProjectsRequest) {
  const res = await service.get<Project[]>("/project/filter", { params: req })
  return Array.isArray(res.data) ? res.data : []
}

export async function createProject(req: CreateProjectRequest) {
  return service.post<Project>("/project", req)
}

export async function updateProject(id: number, req: UpdateProjectRequest) {
  return service.put<Project>(`/project/${id}`, req)
}

export async function deleteProject(id: number) {
  return service.delete(`/project/${id}`)
}
