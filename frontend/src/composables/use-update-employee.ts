import { updateEmployee } from "@/api/employees"
import type { UpdateEmployeeRequest } from "@/model/employees"
import { useMutation, useQueryClient } from "@tanstack/vue-query"
import { toast } from "vue-sonner"

export default function useUpdateEmployee({
  onSuccuss,
  onError
}: {
  onSuccuss?: () => void
  onError?: () => void
}) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, req }: { id: number; req: UpdateEmployeeRequest }) =>
      updateEmployee(id, req),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["employees", "overview"] })
      if (onSuccuss) {
        onSuccuss()
      }
    },
    onError: () => {
      toast.error("更新失败")
      if (onError) {
        onError()
      }
    }
  })
}
