export const ProcessStatusOptions = [
  "等待",
  "准备",
  "进行",
  "返工",
  "完成",
  "暂停",
  "终止"
] as const

export type ProcessStatus = (typeof ProcessStatusOptions)[number]

interface InsertProcessRequest {
  name: string
  projectID: number
  employeeID: number
  description: string
  remark: string
}

interface UpdateProcessRequest {
  name?: string
  description?: string
  remark?: string
  status?: ProcessStatus
  employeeID?: number
  spendTime?: number
  startedAt?: string
  finishedAt?: string
}

interface FinishProcessRequest {
  spendTime: number
  remark: string
}
interface ReworkProcessRequest {
  spendTime: number
  remark: string
}

interface Process {
  id: number
  projectID: number
  name: string
  description: string
  remark: string
  status: ProcessStatus
  employeeID: number
  sequence: number
  spendTime?: number
  startedAt?: string
  finishedAt?: string
}

interface Task {
  id: number
  name: string
  description: string
  projectID: number
  projectName: string
  remainingDays?: number
}

export type {
  InsertProcessRequest,
  UpdateProcessRequest,
  FinishProcessRequest,
  ReworkProcessRequest,
  Process,
  Task
}
