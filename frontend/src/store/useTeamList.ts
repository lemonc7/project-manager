import { create } from "zustand"
import type { Team } from "@/model/teams"

interface TeamList {
  teamList: Team[]
  setTeamList: (teamList: Team[]) => void
}

export const useTeamList = create<TeamList>((set) => ({
  teamList: [],
  setTeamList: (teamList) => set({ teamList })
}))
