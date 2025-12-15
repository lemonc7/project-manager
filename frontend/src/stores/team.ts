import { defineStore } from "pinia"
import { computed, ref } from "vue"
import type { Team } from "../models/teams"

export const useTeamStore = defineStore("team", () => {
  const teamMap = ref(new Map<number, string>())
  function setTeams(list: Team[]) {
    const map = new Map<number, string>()
    list.forEach((t) => map.set(t.id, t.name))
    teamMap.value = map
  }

  function getTeamName(id?: number) {
    if (id === undefined) return "未分配"
    return teamMap.value.get(id) ?? "未知"
  }

  const teamItems = computed(() =>
    Array.from(teamMap.value.entries(), ([id, name]) => ({
      value: id,
      label: name
    }))
  )

  return {
    teamMap,
    teamItems,
    getTeamName,
    setTeams
  }
})
