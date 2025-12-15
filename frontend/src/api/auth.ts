import type {
  ChangePasswordRequest,
  LoginRequest,
  LoginResponse
} from "../models/auth"
import service from "./request"
import { setToken } from "../utils/auth"

export async function login(req: LoginRequest) {
  const res = await service.post<LoginResponse>("/auth/login", req)
  if (!res.data) {
    throw new Error("登录响应数据为空")
  }

  const data = res.data
  if (!data.accessToken || data.userID === undefined || !data.username) {
    console.error("后端响应数据不完整")
    throw new Error("后端响应数据不完整")
  }

  setToken(
    res.data.accessToken,
    res.data.userID,
    res.data.username,
    res.data.roles || []
  )
  return res.data
}

export async function changePassword(req: ChangePasswordRequest) {
  return service.put("/password", req)
}

export async function resetPassword(userID: number, password: string) {
  return service.post(`/password/${userID}`, { password })
}
