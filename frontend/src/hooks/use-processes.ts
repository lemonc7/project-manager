import { listProcesses } from "@/api/processes"
import { useQuery } from "@tanstack/react-query"
import { useParams } from "react-router"

export default function useProcesses() {
  const { id } = useParams<{ id: string }>()
  return useQuery({
    queryKey: ["processes", id],
    queryFn: () => {
      if (!id) {
        throw new Error("缺少流程ID")
      }
      return listProcesses(+id)
    }
  })
}
