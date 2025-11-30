import type { LoginRequest, LoginResponse } from "@/models/auth"
import service from "./request"
import { setToken } from "@/util/auth"

export async function login(req: LoginRequest) {
  const res = await service.post<LoginResponse>("/auth/login", req)
  setToken(
    res.data.accessToken,
    res.data.userID,
    res.data.username,
    res.data.roles || []
  )
  return res.data
}
