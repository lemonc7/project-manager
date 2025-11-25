import { listEmployees } from "@/api/employees"
import type { Employee } from "@/model/employees"
import { useQuery } from "@tanstack/react-query"

export default function useEmployees() {
  return useQuery<Employee[]>({
    queryKey: ["employees", "overview"],
    queryFn: listEmployees
  })
}
