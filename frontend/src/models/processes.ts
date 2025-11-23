type ProcessStatus =
  | "等待"
  | "准备"
  | "进行"
  | "完成"
  | "返工"
  | "终止"
  | "取消"
  | "暂停"

interface InsertProcessRequest {
  name: string
  projectID: number
  employeeID: number
  description?: string
  remark?: string
}

interface UpdateProcessRequest {
  name?: string
  description?: string
  status?: ProcessStatus
  employeeID?: number
  remark?: string
}

interface ListProcessesRequest {
  projectID?: number
  status?: ProcessStatus
  employeeID?: number
}

interface Process {
  id: number
  name: string
  sequence: number
  status: ProcessStatus
  projectName: string
  employeeName: string
  description: string
  remark: string
  startTime: string | null
  endTime: string | null
  createdAt: string
  updatedAt: string
}

interface Task {
  id: number
  name: string
  startTime?: string
  duration?: number
  projectId: number
  projectName: string
}

export type {
  ProcessStatus,
  InsertProcessRequest,
  UpdateProcessRequest,
  ListProcessesRequest,
  Process,
  Task
}
