type ProjectStatus = "排单" | "进行" | "暂停" | "完成" | "终止"
type TaskStatus = "todo" | "ready" | "doing"

interface CreateProjectRequest {
  name: string
  description?: string
  managerID: number
  remark?: string
  duration?: number
}

interface UpdateProjectRequest {
  name?: string
  description?: string
  remark?: string
  status?: ProjectStatus
  managerID?: number
  startedAt?: string
  duration?: number
  finishedAt?: string
}

interface ListProjectsRequest {
  status?: ProjectStatus
  managerID?: number
  startAt?: string
  endAt?: string
}

interface Project {
  id: number
  name: string
  description: string | null
  remark: string | null
  status: ProjectStatus
  managerId: number
  startedAt: string | null
  duration: number | null
  finishedAt: string | null
}

export type {
  ProjectStatus,
  TaskStatus,
  CreateProjectRequest,
  UpdateProjectRequest,
  ListProjectsRequest,
  Project
}
