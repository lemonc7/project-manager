interface CreateEmployeeRequest {
  name: string
  teamID?: number
}

interface UPdateEmployeeRequest {
  id: number
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

export type { CreateEmployeeRequest, UPdateEmployeeRequest, Employee }
