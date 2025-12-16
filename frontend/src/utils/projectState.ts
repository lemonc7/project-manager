import type { ListProjectsRequest } from "../models/projects"

const PROJECT_STATE_KEY = "project-state"

export function getDefaultProjectState(): ListProjectsRequest {
  return {
    managerID: 0,
    status: ["排单", "进行", "暂停"],
    startAt: undefined,
    endAt: undefined
  }
}

export function getProjectState(): ListProjectsRequest {
  const defaultState = getDefaultProjectState()
  const raw = localStorage.getItem(PROJECT_STATE_KEY)
  if (!raw) {
    localStorage.setItem(PROJECT_STATE_KEY, JSON.stringify(defaultState))
    return defaultState
  }

  try {
    const parsed = JSON.parse(raw)
    if (
      parsed.status &&
      Array.isArray(parsed.status) &&
      parsed.status.length === 0
    ) {
      parsed.status = ["排单", "进行", "暂停"]
    }
    localStorage.setItem(PROJECT_STATE_KEY, JSON.stringify(parsed))
    return parsed
  } catch {
    localStorage.setItem(PROJECT_STATE_KEY, JSON.stringify(defaultState))
    return defaultState
  }
}

export function setProjectState(newState: ListProjectsRequest) {
  localStorage.setItem(PROJECT_STATE_KEY, JSON.stringify(newState))
}

export function resetProjectState() {
  localStorage.setItem(
    PROJECT_STATE_KEY,
    JSON.stringify(getDefaultProjectState())
  )
}
