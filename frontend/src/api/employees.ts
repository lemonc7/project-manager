import type {
  CreateEmployeeRequest,
  Employee,
  EmployeeMe,
  UpdateEmployeeRequest
} from "@/model/employees"
import service from "./request"

export async function listEmployees() {
  const res = await service.get<Employee[]>("/employee")
  return Array.isArray(res.data) ? res.data : []
}

export async function getEmployeeMe() {
  const res = await service.get<EmployeeMe>("/employee/me")
  return res.data
}

export async function createEmployee(req: CreateEmployeeRequest) {
  const res = await service.post<Employee>("/employee", req)
  return res.data
}

export async function updateEmployee(id: number, req: UpdateEmployeeRequest) {
  const res = await service.put<Employee>(`/employee/${id}`, req)
  return res.data
}

export async function deleteEmployee(id: number) {
  return service.delete(`/employee/${id}`)
}
