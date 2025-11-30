import type { Role } from "@/models/employees"

const TOKEN_KEY = "access_token"
const USER_ID_KEY = "user_id"
const USERNAME_KEY = "username"
const ROLES_KEY = "roles"

export function setToken(
  token: string,
  userID: number,
  username: string,
  roles: string[]
) {
  localStorage.setItem(TOKEN_KEY, token)
  localStorage.setItem(USER_ID_KEY, userID.toString())
  localStorage.setItem(USERNAME_KEY, username)
  localStorage.setItem(ROLES_KEY, JSON.stringify(roles))
}

export function getToken() {
  return localStorage.getItem(TOKEN_KEY)
}

export function getUserID() {
  const id = localStorage.getItem(USER_ID_KEY)
  return id ? Number(id) : undefined
}

export function removeToken() {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(USER_ID_KEY)
  localStorage.removeItem(USERNAME_KEY)
  localStorage.removeItem(ROLES_KEY)
}

export function isLoggedIn() {
  return !!getToken()
}

export function hasRole(...targetRoles: Role[]) {
  const roles = localStorage.getItem(ROLES_KEY)
  if (!roles) {
    return false
  }
  const parsed = JSON.parse(roles) as Role[]
  return targetRoles.some((role) => parsed.includes(role))
}
