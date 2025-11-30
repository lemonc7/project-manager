import { createEmployee } from "@/api/employees"
import type { CreateEmployeeRequest } from "@/models/employees"
import { useMutation, useQueryClient } from "@tanstack/vue-query"

export default function useCreateEmployee({
  onSuccess,
  onError
}: {
  onSuccess?: () => void
  onError?: () => void
}) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (req: CreateEmployeeRequest) => createEmployee(req),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["employees", "overview"] })
      if (onSuccess) {
        onSuccess()
      }
    },
    onError: () => {
      if (onError) {
        onError()
      }
    }
  })
}
