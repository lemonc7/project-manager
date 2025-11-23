interface CreateEmployeeRequest {
  name: string
  teamID?: number
}

interface UpdateEmployeeRequest {
  name?: string
  teamID?: number
  isActive?: boolean
}

interface Employee {
  id: number
  employeeName: string
  teamName: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export type { CreateEmployeeRequest, UpdateEmployeeRequest, Employee }
