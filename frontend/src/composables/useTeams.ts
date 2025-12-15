import { useMutation, useQuery, useQueryClient } from "@tanstack/vue-query"
import { createTeam, deleteTeam, listTeams, updateTeam } from "../api/teams"
import type { CreateTeamRequest, UpdateTeamRequest } from "../models/teams"

export function useTeams() {
  const queryKey = ["teams", "overview"]
  const queryClient = useQueryClient()

  function invalidateTeams() {
    queryClient.invalidateQueries({ queryKey })
  }

  // 查询数据
  const teamQuery = useQuery({
    queryKey,
    queryFn: listTeams
  })

  // actions
  const createMutation = useMutation({
    mutationFn: (req: CreateTeamRequest) => createTeam(req),
    onSuccess: () => invalidateTeams()
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, req }: { id: number; req: UpdateTeamRequest }) =>
      updateTeam(id, req),
    onSuccess: () => invalidateTeams()
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => deleteTeam(id),
    onSuccess: () => invalidateTeams()
  })

  return {
    teams: teamQuery.data,
    isLoading: teamQuery.isLoading,
    isError: teamQuery.isError,

    createTeam: createMutation.mutateAsync,
    updateTeam: updateMutation.mutateAsync,
    deleteTeam: deleteMutation.mutateAsync,

    isCreating: createMutation.isPending,
    isUpdating: updateMutation.isPending,
    isDeleting: deleteMutation.isPending
  }
}
