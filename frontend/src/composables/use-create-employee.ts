import { createEmployee } from "@/api/employees"
import type { CreateEmployeeRequest } from "@/model/employees"
import { useMutation, useQueryClient } from "@tanstack/vue-query"
import { toast } from "vue-sonner"

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
      toast.success("创建成功")
      if (onSuccess) {
        onSuccess()
      }
    },
    onError: () => {
      toast.error("创建失败")
      if (onError) {
        onError()
      }
    }
  })
}
