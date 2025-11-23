import { isLoggedIn } from "@/util/auth"
import type { JSX } from "react"
import { Navigate } from "react-router"

export default function ProtectElement({
  children
}: {
  children: JSX.Element
}) {
  return isLoggedIn() ? children : <Navigate to="/login" />
}
