import type { CreateTeamRequest, Team, UpdateTeamRequest } from "@/model/teams"
import service from "./request"

export async function listTeams() {
  const res = await service.get<Team[]>("/team")
  return Array.isArray(res.data) ? res.data : []
}

export async function createTeam(req: CreateTeamRequest) {
  const res = await service.post<Team>("/team", req)
  return res.data
}

export async function updateTeam(id: number, req: UpdateTeamRequest) {
  const res = await service.put<Team>(`/team/${id}`, req)
  return res.data
}

export async function deleteTeam(id: number) {
  return service.delete(`/team/${id}`)
}
