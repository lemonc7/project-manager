import { createProject } from "../api/projects"
import type { CreateProjectRequest } from "../models/projects"
import { useMutation, useQueryClient } from "@tanstack/vue-query"

export default function useCreateProject({
  onSuccess,
  onError
}: {
  onSuccess?: () => void
  onError?: () => void
}) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (req: CreateProjectRequest) => createProject(req),
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
