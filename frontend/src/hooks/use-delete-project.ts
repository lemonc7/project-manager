import { deleteProject } from "@/api/projects"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"

export default function useDeleteProject(
  onSuccess?: () => void,
  onError?: () => void
) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => deleteProject(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects", "overview"] })
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
