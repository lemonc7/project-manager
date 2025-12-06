import { roleItems, type Role } from "../models/employees"

export default function getRoleLabels(roles?: Role[]): string {
  return (
    roles
      ?.map((r) => roleItems.find((i) => i.value === r)?.label ?? r)
      .join("，") ?? ""
  )
}
