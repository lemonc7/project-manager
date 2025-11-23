interface CreateTeamRequest {
  name: string
  description?: string
}

interface UpdateTeamRequest {
  name?: string
  description?: string
}

interface Team {
  id: number
  name: string
  description: string | null
}

export type { CreateTeamRequest, UpdateTeamRequest, Team }
