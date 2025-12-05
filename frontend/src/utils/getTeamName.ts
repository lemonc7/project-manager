import type { Team } from "../models/teams"

export default function getTeamName(teamID?: number, teams?: Team[]) {
  if (teamID === undefined || teams === undefined) return "未分组"
  const team = teams.find((t) => t.id === teamID)
  return team ? team.name : "未知小组"
}
