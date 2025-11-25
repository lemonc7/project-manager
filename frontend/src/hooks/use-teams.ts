import { listTeams } from "@/api/teams"
import type { Team } from "@/model/teams"
import { useQuery } from "@tanstack/react-query"

export default function useTeams() {
  return useQuery<Team[]>({
    queryKey: ["teams", "overview"],
    queryFn: listTeams
  })
}
