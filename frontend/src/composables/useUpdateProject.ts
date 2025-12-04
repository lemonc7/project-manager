import { updateProject } from "../api/projects"
import type { UpdateProjectRequest } from "../models/projects"
import { useMutation, useQueryClient } from "@tanstack/vue-query"

export default function useUpdateProject({
  onSuccess,
  onError
}: {
  onSuccess?: () => void
  onError?: () => void
}) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, req }: { id: number; req: UpdateProjectRequest }) =>
      updateProject(id, req),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects", "overview"] })
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
