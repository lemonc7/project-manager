import { useMutation, useQuery, useQueryClient } from "@tanstack/vue-query"
import { defineStore } from "pinia"
import { computed, ref, watch } from "vue"
import { createTeam, deleteTeam, listTeams, updateTeam } from "../api/teams"
import type { CreateTeamRequest } from "../models/teams"

export const useTeamStore = defineStore("team", () => {
  const queryClient = useQueryClient()

  const teamMap = ref<Map<number, string>>(new Map())
  const queryKey = ["teams", "overview"]

  // 查询数据
  const teamQuery = useQuery({
    queryKey,
    queryFn: listTeams
  })

  // 监听数据变化，更新teamMap
  watch(
    teamQuery.data,
    (list) => {
      const newMap = new Map<number, string>()
      list?.forEach((t) => newMap.set(t.id, t.name))
      teamMap.value = newMap
    },
    { immediate: true }
  )

  // 查询teamName
  function getTeamName(id?: number) {
    if (id === undefined) return "未分配"
    return teamMap.value.get(id) ?? "未知"
  }

  // teamItems,用于select
  const teamItems = computed(() =>
    Array.from(teamMap.value.entries(), ([id, name]) => ({
      value: id,
      label: name
    }))
  )

  // actions
  const createMutation = useMutation({
    mutationFn: (req: CreateTeamRequest) => createTeam(req),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey })
    }
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, req }: { id: number; req: CreateTeamRequest }) =>
      updateTeam(id, req),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey })
    }
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => deleteTeam(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey })
    }
  })

  return {
    teams: teamQuery.data,
    isLoading: teamQuery.isLoading,
    isError: teamQuery.isError,

    teamMap,
    teamItems,
    isCreating: createMutation.isPending,
    isUpdating: updateMutation.isPending,
    isDeleting: deleteMutation.isPending,

    getTeamName,

    createTeam: createMutation.mutateAsync,
    updateTeam: updateMutation.mutateAsync,
    deleteTeam: deleteMutation.mutateAsync
  }
})
