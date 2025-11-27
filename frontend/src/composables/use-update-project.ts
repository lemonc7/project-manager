import { updateProject } from "@/api/projects"
import type { UpdateProjectRequest } from "@/model/projects"
import { useMutation, useQueryClient } from "@tanstack/vue-query"
import { toast } from "vue-sonner"

export default function useUpdateProject({
  onSuccuss,
  onError
}: {
  onSuccuss?: () => void
  onError?: () => void
}) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, req }: { id: number; req: UpdateProjectRequest }) =>
      updateProject(id, req),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects", "overview"] })
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
