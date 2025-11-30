export const RoleMap = {
  admin: "管理员",
  manager: "项目经理",
  employee: "员工"
} as const

// 自动推导类型为 "admin" | "manager" | "employee"
export type Role = keyof typeof RoleMap

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
  roles?: string[]
  teamID?: number
}

export type { CreateEmployeeRequest, UpdateEmployeeRequest, Employee }
