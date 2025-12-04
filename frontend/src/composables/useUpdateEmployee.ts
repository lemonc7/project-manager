import { updateEmployee } from "../api/employees"
import type { UpdateEmployeeRequest } from "../models/employees"
import { useMutation, useQueryClient } from "@tanstack/vue-query"

export default function useUpdateEmployee({
  onSuccess,
  onError
}: {
  onSuccess?: () => void
  onError?: () => void
}) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, req }: { id: number; req: UpdateEmployeeRequest }) =>
      updateEmployee(id, req),
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
