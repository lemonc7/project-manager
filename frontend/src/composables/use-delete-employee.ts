import { deleteEmployee } from "@/api/employees"
import { useMutation, useQueryClient } from "@tanstack/vue-query"
import { toast } from "vue-sonner"

export default function useDeleteEmployee({
  onSuccess,
  onError
}: {
  onSuccess?: () => void
  onError?: () => void
}) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => deleteEmployee(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["employees", "overview"] })
      toast.success("删除成功")
      if (onSuccess) {
        onSuccess()
      }
    },
    onError: () => {
      toast.error("删除失败")
      if (onError) {
        onError()
      }
    }
  })
}
