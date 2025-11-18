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

interface InsertProcessRequest {
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
  startTime: NullableTime
  endTime: NullableTime
  createdAt: string
  updatedAt: string
}

export type {
  ProcessStatus,
  InsertProcessRequest,
  UpdateProcessRequest,
  ListProcessesRequest,
  Process
}
