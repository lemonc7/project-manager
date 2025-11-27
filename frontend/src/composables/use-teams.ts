import { listTeams } from "@/api/teams"
import { useQuery } from "@tanstack/vue-query"

export default function useTeams() {
  return useQuery({
    queryKey: ["teams", "overview"],
    queryFn: listTeams
  })
}
