import { createProject } from "@/api/projects"
import type { CreateProjectRequest } from "@/model/projects"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"

export default function useCreateProject(
  onSuccuss?: () => void,
  onError?: () => void
) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (req: CreateProjectRequest) => createProject(req),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects", "overview"] })
      toast.success("创建成功")
      if (onSuccuss) {
        onSuccuss()
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
