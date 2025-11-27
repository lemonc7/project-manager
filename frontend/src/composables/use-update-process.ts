import { updateProcess } from "@/api/processes"
import type { UpdateProcessRequest } from "@/model/processes"
import { useMutation, useQueryClient } from "@tanstack/vue-query"
import { useRoute } from "vue-router"
import { toast } from "vue-sonner"

export default function useUpdateProcess({
  onSuccess,
  onError
}: {
  onSuccess?: () => void
  onError?: () => void
}) {
  const queryClient = useQueryClient()
  const route = useRoute()
  const id = route.params.id as string | undefined
  return useMutation({
    mutationFn: ({ id, req }: { id: number; req: UpdateProcessRequest }) =>
      updateProcess(id, req),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["processes", id] })
      if (onSuccess) {
        onSuccess()
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
