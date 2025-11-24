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
  startedAt?: string
  finishedAt?: string
}

interface Process {
  id: number
  projectId: number
  name: string
  description: string | null
  remark: string | null
  status: ProcessStatus
  employeeId: number | null
  sequence: number
  startedAt: string | null
  finishedAt: string | null
}

interface Task {
  id: number
  name: string
  description?: string
  startedAt?: string
  remainingDays?: number
  projectId: number
  projectName: string
}

export type { InsertProcessRequest, UpdateProcessRequest, Process, Task }
