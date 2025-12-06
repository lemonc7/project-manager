export const roleItems = [
  { value: "admin", label: "管理员" },
  { value: "manager", label: "项目经理" },
  { value: "employee", label: "员工" },
] as const

export type Role = typeof roleItems[number]['value']
export type RoleItem = typeof roleItems[number]


interface CreateEmployeeRequest {
  name: string
  username: string
  password: string
  roles?: Role[]
  teamID?: number
}

interface UpdateEmployeeRequest {
  name?: string
  roles?: Role[]
  teamID?: number
  isActive?: boolean
}

interface Employee {
  id: number
  name: string
  roles?: Role[]
  teamID?: number
}

export type { CreateEmployeeRequest, UpdateEmployeeRequest, Employee }
