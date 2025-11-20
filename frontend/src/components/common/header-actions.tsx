import { useEffect, useState } from "react"
import { createPortal } from "react-dom"

export const HEADER_ACTIONS_ID = "header-action-slot"

export function HeaderActions({ children }: { children: React.ReactNode }) {
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
    return () => setMounted(false)
  }, [])

  const container = document.getElementById(HEADER_ACTIONS_ID)

  if (!mounted || !container) return null

  return createPortal(children, container)
}
