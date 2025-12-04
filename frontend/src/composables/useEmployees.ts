import { listEmployees } from "../api/employees"
import { useQuery } from "@tanstack/vue-query"

export default function useEmployees() {
  return useQuery({
    queryKey: ["employees", "overview"],
    queryFn: listEmployees
  })
}
