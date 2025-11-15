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
  return service.post("/employee", req)
}

export async function updateEmployee(
  id: number,
  req: Omit<UPdateEmployeeRequest, "id">
) {
  return service.put(`/employee/${id}`, req)
}

export async function deleteEmployee(id: number) {
  return service.delete(`/employee/${id}`)
}
