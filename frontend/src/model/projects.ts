type ProjectStatus = "排单" | "进行" | "暂停" | "完成" | "终止" | "归档"

interface CreateProjectRequest {
  name: string
  description: string
  managerID: number
  remark: string
}

interface UpdateProjectRequest {
  id: number
  name?: string
  description?: string
  status?: ProjectStatus
  managerID?: number
  startTime?: Date
  duration?: number
  remark?: string
}

interface ListProjectsRequest {
  status?: ProjectStatus
  managerID?: number
  startTime?: Date
  endTime?: Date
}

interface NullableTime {
  Time: string
  Valid: boolean
}

interface Project {
  id: number
  name: string
  managerName: string
  description: string
  remark: string
  status: ProjectStatus
  startTime: NullableTime
  duration: number
  createdAt: string
  updatedAt: string
}

export type {
  ProjectStatus,
  CreateProjectRequest,
  UpdateProjectRequest,
  ListProjectsRequest,
  NullableTime,
  Project
}
