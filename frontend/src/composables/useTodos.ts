import { useMutation, useQuery, useQueryClient } from "@tanstack/vue-query"
import { getUserID } from "../utils/auth"
import {
  createTodo,
  deleteTodo,
  listTodosByReporter,
  updateTodo
} from "../api/todos"
import type {
  CreateTodoRequest,
  TodoStatus,
  UpdateTodoRequest
} from "../models/todos"
import { getTodoState } from "../utils/todoState"
import { computed, ref } from "vue"

function buildTodoKey(
  id: number,
  todoType: "reporter" | "assignee",
  status?: TodoStatus[]
) {
  return ["todos", id, todoType, status ?? null]
}

export function useTodos() {
  const userID = getUserID()
  const queryClient = useQueryClient()

  const initialState = getTodoState()
  const req = ref({
    reporterStatus: initialState.reporterStatus,
    assigneeStatus: initialState.assigneeStatus
  })
  const reporterQueryKey = computed(() => {
    if (!userID) {
      throw new Error("缺少用户ID")
    }
    return buildTodoKey(userID, "reporter", req.value.reporterStatus)
  })

  const assigneeQueryKey = computed(() => {
    if (!userID) {
      throw new Error("缺少用户ID")
    }
    return buildTodoKey(userID, "assignee", req.value.assigneeStatus)
  })

  function invalidateTodos() {
    queryClient.invalidateQueries({
      queryKey: ["todos", userID]
    })
  }

  const todosReporterQuery = useQuery({
    queryKey: reporterQueryKey,
    queryFn: () => {
      if (!userID) {
        throw new Error("缺少用户ID")
      }
      return listTodosByReporter(userID, req.value.reporterStatus)
    },
    refetchInterval: 1000 * 60,
    staleTime: Infinity,
    initialData: []
  })
  
  const todosAssigneeQuery = useQuery({
    queryKey: assigneeQueryKey,
    queryFn: () => {
      if (!userID) {
        throw new Error("缺少用户ID")
      }
      return listTodosByReporter(userID, req.value.assigneeStatus)
    },
    refetchInterval: 1000 * 60,
    staleTime: Infinity,
    initialData: []
  })

  // actions
  const createMutation = useMutation({
    mutationFn: (req: CreateTodoRequest) => createTodo(req),
    onSuccess: () => invalidateTodos()
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, req }: { id: number; req: UpdateTodoRequest }) =>
      updateTodo(id, req),
    onSuccess: () => invalidateTodos()
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => deleteTodo(id),
    onSuccess: () => invalidateTodos()
  })

  return {
    todosReporter: todosReporterQuery.data,
    todosAssignee: todosAssigneeQuery.data,
    isReporterLoading: todosReporterQuery.isLoading,
    isAssigneeLoading: todosAssigneeQuery.isLoading,
    isReporterError: todosReporterQuery.isError,
    isAssigneeError: todosAssigneeQuery.isError,

    isCreating: createMutation.isPending,
    isUpdating: updateMutation.isPending,
    isDeleting: deleteMutation.isPending,

    setQuery: (newState: {
      reporterStatus: TodoStatus[]
      assigneeStatus: TodoStatus[]
    }) => {
      req.value = newState
    },

    createTodo: createMutation.mutateAsync,
    updateTodo: updateMutation.mutateAsync,
    deleteTodo: deleteMutation.mutateAsync
  }
}
