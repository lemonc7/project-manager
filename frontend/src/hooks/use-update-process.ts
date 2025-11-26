import { updateProcess } from "@/api/processes"
import type { UpdateProcessRequest } from "@/model/processes"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { useParams } from "react-router"
import { toast } from "sonner"

export default function useUpdateProcess(
  onSuccess?: () => void,
  onError?: () => void
) {
  const queryClient = useQueryClient()
  const { id } = useParams<{ id: string }>()
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
