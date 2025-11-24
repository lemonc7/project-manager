interface CreateEmployeeRequest {
  name: string
  username: string
  password: string
  roles?: string[]
  teamID?: number
}

interface UpdateEmployeeRequest {
  name?: string
  roles?: string[]
  teamID?: number
  isActive?: boolean
}

interface Employee {
  id: number
  name: string
  roles: string[]
  teamId: number | null
}

interface EmployeeMe {
  id: number
  name: string
  roles: string[]
}

export type {
  CreateEmployeeRequest,
  UpdateEmployeeRequest,
  Employee,
  EmployeeMe
}
