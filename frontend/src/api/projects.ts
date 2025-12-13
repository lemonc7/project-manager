import type {
  CopyProjectRequest,
  CreateProjectRequest,
  ListProjectsRequest,
  Project,
  ProjectGet,
  UpdateProjectRequest
} from "../models/projects"
import service from "./request"

export async function listProjects(req: ListProjectsRequest) {
  const res = await service.get<Project[]>(`/project`, {
    params: req
  })
  return Array.isArray(res.data) ? res.data : []
}

export async function createProject(req: CreateProjectRequest) {
  return service.post<Project>("/project", req)
}

export async function updateProject(id: number, req: UpdateProjectRequest) {
  return service.patch<Project>(`/project/${id}`, req)
}

export async function deleteProject(id: number) {
  return service.delete(`/project/${id}`)
}

export async function startProject(id: number) {
  return service.post(`/project/start/${id}`)
}

export async function finishProject(id: number) {
  return updateProject(id, { status: "完成" })
}

export async function getProject(id: number) {
  const res = await service.get<ProjectGet>(`/project/${id}`)
  return res.data
}

export async function resetProject(id: number) {
  return service.post(`/project/reset/${id}`)
}

export async function copyProject(id: number, req: CopyProjectRequest) {
  return service.post(`/project/copy/${id}`, req)
}

export async function isProjectExists(name: string) {
  const res = await service.get<boolean>("/project/exist", {
    params: {
      newName: name
    }
  })
  return res.data
}
