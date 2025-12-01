import { finishProcess } from "@/api/processes"
import type { FinishProcessRequest } from "@/models/processes"
import { useMutation, useQueryClient } from "@tanstack/vue-query"
import { useRoute } from "vue-router"

export default function useFinishProcess({
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
    mutationFn: ({ id, req }: { id: number; req: FinishProcessRequest }) => {
      return finishProcess(id, req)
    },
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
