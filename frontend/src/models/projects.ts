export type TaskStatus = "todo" | "ready" | "doing"

export const ProjectStatusOptions = [
  "排单",
  "进行",
  "暂停",
  "完成",
  "终止"
] as const

export type ProjectStatus = (typeof ProjectStatusOptions)[number]

interface CreateProjectRequest {
  name: string
  description: string
  managerID: number
  remark: string
  priority: number
  duration?: number
}

interface UpdateProjectRequest {
  name?: string
  description?: string
  remark?: string
  status?: ProjectStatus
  managerID?: number
  priority?: number
  duration?: number
}

interface ListProjectsRequest {
  startAt?: string
  endAt?: string
}

interface Project {
  id: number
  name: string
  description: string
  remark: string
  status: ProjectStatus
  managerID: number
  priority: number
  startedAt?: string
  duration?: number
  finishedAt?: string
  remainingDays?: number
}

interface ProjectGet {
  name: string
  description: string
  remainingDays?: number
}

export type {
  CreateProjectRequest,
  UpdateProjectRequest,
  ListProjectsRequest,
  Project,
  ProjectGet
}
