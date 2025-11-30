import { insertProcessAt, insertProcessLast } from "@/api/processes"
import type { InsertProcessRequest } from "@/models/processes"
import { useMutation, useQueryClient } from "@tanstack/vue-query"
import { useRoute } from "vue-router"

export default function useCreateProcess({
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
    mutationFn: ({
      sequence,
      req
    }: {
      sequence?: number
      req: Omit<InsertProcessRequest, "projectID">
    }) => {
      if (sequence) {
        return insertProcessAt(sequence, { ...req, projectID })
      } else {
        return insertProcessLast({ ...req, projectID })
      }
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
