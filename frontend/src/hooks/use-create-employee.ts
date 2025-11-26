import { createEmployee } from "@/api/employees"
import type { CreateEmployeeRequest } from "@/model/employees"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"

export default function useCreateEmployee(
  onSuccess?: () => void,
  onError?: () => void
) {
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
