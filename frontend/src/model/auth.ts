interface LoginRequest {
  username: string
  password: string
}

interface LoginResponse {
  accessToken: string
  userID: number
  username: string
  roles?: string[]
}

export type { LoginRequest, LoginResponse }
