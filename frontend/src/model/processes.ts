import type { NullableTime } from "./projects"

type ProcessStatus =
  | "等待"
  | "准备"
  | "进行"
  | "完成"
  | "返工"
  | "终止"
  | "取消"
  | "暂停"

interface CreateProcessRequest {
  name: string
  projectID: number
  employeeID: number
  description: string
  sequence: number
  remark: string
}

interface UpdateProcessRequest {
  id: number
  name?: string
  description?: string
  status?: ProcessStatus
  employeeID?: number
  sequence?: number
  remark?: string
}

interface ListProcessesRequest {
  projectID?: number
  status?: ProcessStatus
  employeeID?: number
}

interface Process {
  ID: number
  processName: string
  sequence: number
  status: ProcessStatus
  projectName: string
  employeeName: string
  description: string
  startTime: NullableTime
  endTime: NullableTime
  createdAt: string
  updatedAt: string
}

export type {
  ProcessStatus,
  CreateProcessRequest,
  UpdateProcessRequest,
  ListProcessesRequest,
  Process
}
