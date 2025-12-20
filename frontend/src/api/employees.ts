import type {
  CreateEmployeeRequest,
  Employee,
  UpdateEmployeeRequest
} from "../models/employees"
import service from "./request"

export async function listEmployees() {
  const res = await service.get<Employee[]>("/employee")
  return res.data
}

export async function getEmployeeName() {
  const res = await service.get<string>("/employee/me")
  return res.data || ""
}

export async function createEmployee(req: CreateEmployeeRequest) {
  return service.post("/employee", req)
}

export async function updateEmployee(id: number, req: UpdateEmployeeRequest) {
  return service.patch(`/employee/${id}`, req)
}

export async function deleteEmployee(id: number) {
  return service.delete(`/employee/${id}`)
}
