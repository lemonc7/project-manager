type ProjectStatus = "排单" | "进行" | "暂停" | "完成" | "终止" | "归档"
type TaskStatus = "todo" | "ready" | "doing"

interface CreateProjectRequest {
  name: string
  description?: string
  managerID: number
  remark?: string
}

interface UpdateProjectRequest {
  name?: string
  description?: string
  status?: ProjectStatus
  managerID?: number
  startTime?: string
  duration?: number
  remark?: string
}

interface ListProjectsRequest {
  status?: ProjectStatus
  managerID?: number
  startTime?: string
  endTime?: string
}

interface Project {
  id: number
  name: string
  managerName: string
  description: string
  remark: string
  status: ProjectStatus
  startTime: string | null
  duration: number
  createdAt: string
  updatedAt: string
}

export type {
  ProjectStatus,
  TaskStatus,
  CreateProjectRequest,
  UpdateProjectRequest,
  ListProjectsRequest,
  Project
}
