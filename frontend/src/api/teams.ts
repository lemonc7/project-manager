import type { CreateTeamRequest, Team, UpdateTeamRequest } from "@/model/teams"
import service from "./request"

export async function listTeams() {
  const res = await service.get<Team[]>("/team")
  return Array.isArray(res.data) ? res.data : []
}

export async function createTeam(req: CreateTeamRequest) {
  return service.post("/team", req)
}

export async function updateTeam(
  id: number,
  req: Omit<UpdateTeamRequest, "id">
) {
  return service.put(`/team/${id}`, req)
}

export async function deleteTeam(id: number) {
  return service.delete(`/team/${id}`)
}
