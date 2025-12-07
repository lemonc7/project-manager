import { roleMap, type Role } from "../models/employees"

export default function getRoleLabels(roles?: Role[]): string {
  return roles?.map((r) => roleMap[r] ?? r).join("，") ?? ""
}
