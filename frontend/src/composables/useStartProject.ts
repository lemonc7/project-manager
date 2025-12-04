import { startProject } from "../api/projects"
import { useMutation, useQueryClient } from "@tanstack/vue-query"

export default function useStartProject({
  onSuccess,
  onError
}: {
  onSuccess?: () => void
  onError?: () => void
}) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => {
      return startProject(id)
    },
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
