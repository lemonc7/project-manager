interface CreateTeamRequest {
  name: string
  description: string
}

interface UpdateTeamRequest {
  name?: string
  description?: string
}

interface Team {
  id: number
  name: string
  description: string
}

export type { CreateTeamRequest, UpdateTeamRequest, Team }
