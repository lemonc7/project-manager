import { useMutation, useQuery, useQueryClient } from "@tanstack/vue-query"
import { defineStore } from "pinia"
import { computed, ref, watch } from "vue"
import { createEmployee, listEmployees, updateEmployee } from "../api/employees"
import type {
  CreateEmployeeRequest,
  UpdateEmployeeRequest
} from "../models/employees"

export const useEmployeeStore = defineStore("employee", () => {
  const queryClient = useQueryClient()

  const employeeMap = ref<Map<number, string>>(new Map())
  const queryKey = ["employees", "overview"]
  // 加载数据
  const employeeQuery = useQuery({
    queryKey,
    queryFn: listEmployees
  })

  // 监听数据变化，更新employeeMap
  watch(
    employeeQuery.data,
    (list) => {
      const newMap = new Map<number, string>()
      list?.forEach((e) => newMap.set(e.id, e.name))
      employeeMap.value = newMap
    },
    { immediate: true }
  )

  // 查询employeeName
  const getEmployeeName = (id?: number) => {
    if (id === undefined) return "未分配"
    return employeeMap.value.get(id) ?? "未知"
  }

  // employeeItems，用于select
  const employeeItems = computed(() =>
    Array.from(employeeMap.value.entries(), ([id, name]) => ({
      value: id,
      label: name
    }))
  )

  // actions
  const createMutation = useMutation({
    mutationFn: (req: CreateEmployeeRequest) => createEmployee(req),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey })
    }
  })
  const updateMutation = useMutation({
    mutationFn: ({ id, req }: { id: number; req: UpdateEmployeeRequest }) =>
      updateEmployee(id, req),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey })
    }
  })
  // 软删除
  const deleteMutation = useMutation({
    mutationFn: (id: number) => updateEmployee(id, { isActive: false }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey })
    }
  })

  return {
    employees: employeeQuery.data,
    isLoading: employeeQuery.isLoading,
    isError: employeeQuery.isError,

    employeeMap,
    employeeItems,

    isCreating: createMutation.isPending,
    isUpdating: updateMutation.isPending,
    isDeleting: deleteMutation.isPending,

    getEmployeeName,

    createEmployee: createMutation.mutateAsync,
    updateEmployee: updateMutation.mutateAsync,
    deleteEmployee: deleteMutation.mutateAsync
  }
})
