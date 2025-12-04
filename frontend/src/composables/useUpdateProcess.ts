import { updateProcess } from "../api/processes"
import type { UpdateProcessRequest } from "../models/processes"
import { useMutation, useQueryClient } from "@tanstack/vue-query"
import { useRoute } from "vue-router"

export default function useUpdateProcess({
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
    mutationFn: ({ id, req }: { id: number; req: UpdateProcessRequest }) =>
      updateProcess(id, req),
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
