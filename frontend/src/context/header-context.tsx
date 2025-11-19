/* eslint-disable react-refresh/only-export-components */
import {
  createContext,
  useCallback,
  useContext,
  useState,
  type ReactNode
} from "react"

export interface HeaderAction {
  key: string
  element: ReactNode
}

interface HeaderContextType {
  actions: HeaderAction[]
  setActions: (actions: HeaderAction[]) => void
}

const HeaderContext = createContext<HeaderContextType | null>(null)

function HeaderProvider({ children }: { children: ReactNode }) {
  const [actions, setActionsState] = useState<HeaderAction[]>([])

  // 避免重复re-render
  const setActions = useCallback((actions: HeaderAction[]) => {
    setActionsState(actions)
  }, [])

  return (
    <HeaderContext.Provider value={{ actions, setActions }}>
      {children}
    </HeaderContext.Provider>
  )
}

function useHeader() {
  const ctx = useContext(HeaderContext)
  if (!ctx) {
    throw new Error("useHeader must be used within a HeaderProvider")
  }
  return ctx
}

export { HeaderProvider, useHeader }
