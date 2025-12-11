import { useMutation, useQuery, useQueryClient } from "@tanstack/vue-query"
import { defineStore } from "pinia"
import { useRoute } from "vue-router"
import {
  deleteProcess,
  finishProcess,
  insertProcessAt,
  insertProcessLast,
  listProcesses,
  reworkProcess,
  updateProcess
} from "../api/processes"
import type {
  FinishProcessRequest,
  InsertProcessRequest,
  ReworkProcessRequest,
  UpdateProcessRequest
} from "../models/processes"
import { computed, ref } from "vue"

export const useProcessStore = defineStore("process", () => {
  const route = useRoute()
  const projectID = computed(() => parseInt(route.params.id as string))
  const queryClient = useQueryClient()
  // 查询数据
  const processQuery = useQuery({
    queryKey: ["processes", projectID],
    queryFn: () => {
      const id = projectID.value
      if (!id || isNaN(id)) {
        throw new Error("缺少流程ID")
      }
      return listProcesses(id)
    },
    enabled: computed(() => !!projectID.value && !isNaN(projectID.value))
  })

  function invalidateProcesses() {
    queryClient.invalidateQueries({
      queryKey: ["processes", projectID]
    })
  }

  // 派生 maxSequence
  const maxSequence = computed(() => {
    const processes = processQuery.data?.value
    if (processes?.length) {
      const maxVal = processes.reduce((max, p) => Math.max(max, p.sequence), 0)
      return maxVal + 1
    }
    return 1
  })

  // 新增状态: 存储要滚动的行序号
  const sequenceToScroll = ref<number | null>(null)

  const setSequenceToScroll = (seq: number) => {
    sequenceToScroll.value = seq
  }

  const clearSequenceToScroll = () => {
    sequenceToScroll.value = null
  }

  // actions
  const createMutation = useMutation({
    mutationFn: ({
      sequence,
      req
    }: {
      sequence?: number
      req: Omit<InsertProcessRequest, "projectID">
    }) => {
      const id = projectID.value
      if (!id || isNaN(id)) {
        throw new Error("缺少项目ID，无法创建流程")
      }
      if (sequence) {
        return insertProcessAt(sequence, { ...req, projectID: id })
      } else {
        return insertProcessLast({ ...req, projectID: id })
      }
    },
    onSuccess: () => invalidateProcesses()
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, req }: { id: number; req: UpdateProcessRequest }) =>
      updateProcess(id, req),
    onSuccess: () => invalidateProcesses()
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => deleteProcess(id),
    onSuccess: () => invalidateProcesses()
  })

  const finishMutation = useMutation({
    mutationFn: ({ id, req }: { id: number; req: FinishProcessRequest }) =>
      finishProcess(id, req),
    onSuccess: () => invalidateProcesses()
  })

  const reworkMutation = useMutation({
    mutationFn: ({ id, req }: { id: number; req: ReworkProcessRequest }) =>
      reworkProcess(id, req),
    onSuccess: () => invalidateProcesses()
  })

  return {
    processes: processQuery.data,
    isLoading: processQuery.isLoading,
    isError: processQuery.isError,

    maxSequence,
    sequenceToScroll,

    isCreating: createMutation.isPending,
    isUpdating: updateMutation.isPending,
    isDeleting: deleteMutation.isPending,
    isFinishing: finishMutation.isPending,
    isReworking: reworkMutation.isPending,

    setSequenceToScroll,
    clearSequenceToScroll,
    createProcess: createMutation.mutateAsync,
    updateProcess: updateMutation.mutateAsync,
    deleteProcess: deleteMutation.mutateAsync,
    finishProcess: finishMutation.mutateAsync,
    reworkProcess: reworkMutation.mutateAsync
  }
})
