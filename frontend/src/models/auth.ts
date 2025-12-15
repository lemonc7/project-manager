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

interface ChangePasswordRequest {
  oldPassword: string
  newPassword: string
}

export type { LoginRequest, LoginResponse, ChangePasswordRequest }
