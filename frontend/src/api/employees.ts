import type {
  CreateEmployeeRequest,
  Employee,
  UpdateEmployeeRequest
} from "@/model/employees"
import service from "./request"

export async function listEmployees() {
  const res = await service.get<Employee[]>("/employee")
  return Array.isArray(res.data) ? res.data : []
}

export async function getEmployeeMe() {
  const res = await service.get<Employee>("/employee/me")
  return res.data
}

export async function createEmployee(req: CreateEmployeeRequest) {
  return service.post<Employee>("/employee", req)
}

export async function updateEmployee(id: number, req: UpdateEmployeeRequest) {
  return service.put<Employee>(`/employee/${id}`, req)
}

export async function deleteEmployee(id: number) {
  return service.delete(`/employee/${id}`)
}
