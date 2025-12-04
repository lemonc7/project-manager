import { deleteProcess } from "../api/processes"
import { useMutation, useQueryClient } from "@tanstack/vue-query"
import { useRoute } from "vue-router"

export default function useDeleteProcess({
  onSuccess,
  onError
}: {
  onSuccess?: () => void
  onError?: () => void
}) {
  const queryClient = useQueryClient()
  const route = useRoute()
  const projectID = parseInt(route.params.id as string)

  return useMutation({
    mutationFn: (id: number) => deleteProcess(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["processes", projectID] })
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
