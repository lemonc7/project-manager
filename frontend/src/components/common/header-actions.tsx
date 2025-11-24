import { createPortal } from "react-dom"

export const HEADER_ACTIONS_ID = "header-action-slot"

export default function HeaderActions({
  children
}: {
  children: React.ReactNode
}) {
  const container = document.getElementById(HEADER_ACTIONS_ID)

  return container ? createPortal(children, container) : null
}
