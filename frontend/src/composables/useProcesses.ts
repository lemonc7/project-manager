import { useQuery } from "@tanstack/vue-query"
import { useRoute } from "vue-router"
import { listProcesses } from "../api/processes"
import { computed } from "vue"

export function useProcesses() {
  const route = useRoute()
  const id = parseInt(route.params.id as string)

  const query = useQuery({
    queryKey: ["processes", id],
    queryFn: () => {
      if (!id) {
        throw new Error("缺少流程ID")
      }
      return listProcesses(id)
    }
  })

  // 派生 maxSequence
  const maxSequence = computed(() => {
    const processes = query.data?.value
    if (processes?.length) {
      const maxVal = processes.reduce((max, p) => Math.max(max, p.sequence), 0)
      return maxVal + 1
    }
    return 1
  })

  return {
    ...query,
    maxSequence
  }
}
