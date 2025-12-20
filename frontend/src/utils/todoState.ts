import type { TodoStatus } from "../models/todos"

const TODO_STATE_KEY = "todo-state"

interface TodoState {
  reporterStatus: TodoStatus[]
  assigneeStatus: TodoStatus[]
}

export function getDefaultTodoState(): TodoState {
  return {
    reporterStatus: ["doing", "todo"],
    assigneeStatus: ["doing", "todo"]
  }
}

export function getTodoState(): TodoState {
  const defaultState = getDefaultTodoState()
  const raw = localStorage.getItem(TODO_STATE_KEY)
  if (!raw) {
    // 本地没存储，用默认值并存储
    localStorage.setItem(TODO_STATE_KEY, JSON.stringify(defaultState))
    return defaultState
  }

  try {
    const parsed = JSON.parse(raw)
    // 本地数据有问题，用默认值并存储
    if (
      parsed.status &&
      Array.isArray(parsed.status) &&
      parsed.status.length === 0
    ) {
      parsed.status = ["doing", "todo"]
    }
    localStorage.setItem(TODO_STATE_KEY, JSON.stringify(parsed))
    return parsed
  } catch {
    // 报错，用默认值并存储
    localStorage.setItem(TODO_STATE_KEY, JSON.stringify(defaultState))
    return defaultState
  }
}

export function setTodoState(newState: { status: TodoStatus[] }) {
  localStorage.setItem(TODO_STATE_KEY, JSON.stringify(newState))
}

export function resetTodoState() {
  localStorage.setItem(TODO_STATE_KEY, JSON.stringify(getDefaultTodoState()))
}
