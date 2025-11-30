import { useQuery } from "@tanstack/vue-query"
import { useRoute } from "vue-router"
import { listProcesses } from "@/api/processes"

export function useProcesses() {
  const route = useRoute()
  const id = parseInt(route.params.id as string)

  return useQuery({
    queryKey: ["processes", id],
    queryFn: () => {
      if (!id) {
        throw new Error("缺少流程ID")
      }
      return listProcesses(id)
    }
  })
}
