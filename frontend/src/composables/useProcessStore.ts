import { useMutation, useQuery, useQueryClient } from "@tanstack/vue-query"
import { defineStore } from "pinia"
import { useRoute } from "vue-router"
import {
  deleteProcess,
  finishProcess,
  insertProcessAt,
  insertProcessLast,
  listProcesses,
  updateProcess
} from "../api/processes"
import type {
  InsertProcessRequest,
  UpdateProcessRequest
} from "../models/processes"
import { computed } from "vue"

export const useProcessStore = defineStore("process", () => {
  const route = useRoute()
  const id = parseInt(route.params.id as string)
  const queryClient = useQueryClient()
  // 查询数据
  const processQuery = useQuery({
    queryKey: ["processes", id],
    queryFn: () => {
      if (!id) {
        throw new Error("缺少流程ID")
      }
      return listProcesses(id)
    }
  })

  // 派生 maxSequence
  const maxSequence = computed(() => {
    const processes = processQuery.data?.value
    if (processes?.length) {
      const maxVal = processes.reduce((max, p) => Math.max(max, p.sequence), 0)
      return maxVal + 1
    }
    return 1
  })

  // actions
  const createMutation = useMutation({
    mutationFn: ({
      sequence,
      req
    }: {
      sequence?: number
      req: Omit<InsertProcessRequest, "projectID">
    }) => {
      if (sequence) {
        return insertProcessAt(sequence, { ...req, projectID: id })
      } else {
        return insertProcessLast({ ...req, projectID: id })
      }
    },
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: ["processes", id] })
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, req }: { id: number; req: UpdateProcessRequest }) =>
      updateProcess(id, req),
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: ["processes", id] })
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => deleteProcess(id),
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: ["processes", id] })
  })

  const finishMutation = useMutation({
    mutationFn: ({ id, spendTime }: { id: number; spendTime: number }) =>
      finishProcess(id, { spendTime }),
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: ["processes", id] })
  })

  return {
    processes: processQuery.data,
    isLoading: processQuery.isLoading,
    isError: processQuery.isError,

    maxSequence,

    isCreating: createMutation.isPending,
    isUpdating: updateMutation.isPending,
    isDeleting: deleteMutation.isPending,
    isFinishing: finishMutation.isPending,

    createProcess: createMutation.mutateAsync,
    updateProcess: updateMutation.mutateAsync,
    deleteProcess: deleteMutation.mutateAsync,
    finishProcess: finishMutation.mutateAsync
  }
})
