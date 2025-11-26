import { listEmployees } from "@/api/employees"
import { useQuery } from "@tanstack/react-query"

export default function useEmployees() {
  return useQuery({
    queryKey: ["employees", "overview"],
    queryFn: listEmployees
  })
}
