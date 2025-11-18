import type {
  CreateEmployeeRequest,
  Employee,
  UPdateEmployeeRequest
} from "@/model/employees"
import service from "./request"

export async function listEmployees() {
  const res = await service.get<Employee[]>("/employee")
  return Array.isArray(res.data) ? res.data : []
}

export async function createEmployee(req: CreateEmployeeRequest) {
  const res = await service.post<Employee>("/employee", req)
  return res.data
}

export async function updateEmployee(req: UPdateEmployeeRequest) {
  const { id, ...rest } = req
  const res = await service.put<Employee>(`/employee/${id}`, rest)
  return res.data
}

export async function deleteEmployee(id: number) {
  return service.delete(`/employee/${id}`)
}
