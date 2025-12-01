export const ProcessStatusOptions = [
  "等待",
  "准备",
  "进行",
  "完成",
  "暂停",
  "终止"
] as const

export type ProcessStatus = (typeof ProcessStatusOptions)[number]

interface InsertProcessRequest {
  name: string
  projectID: number
  employeeID?: number
  description?: string
  remark?: string
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
}

interface Process {
  id: number
  projectID: number
  name: string
  description?: string
  remark?: string
  status: ProcessStatus
  employeeID?: number
  spendTime?: number
  sequence: number
  startedAt?: string
  finishedAt?: string
}

interface Task {
  id: number
  name: string
  description?: string
  startedAt?: string
  remainingDays?: number
  projectID: number
  projectName: string
}

export type {
  InsertProcessRequest,
  UpdateProcessRequest,
  FinishProcessRequest,
  Process,
  Task
}
