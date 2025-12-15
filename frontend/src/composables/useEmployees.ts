import { useMutation, useQuery, useQueryClient } from "@tanstack/vue-query"
import { createEmployee, listEmployees, updateEmployee } from "../api/employees"
import type {
  CreateEmployeeRequest,
  UpdateEmployeeRequest
} from "../models/employees"

export function useEmployees() {
  const queryKey = ["employees", "overview"]
  const queryClient = useQueryClient()

  // 加载数据
  const employeeQuery = useQuery({
    queryKey,
    queryFn: listEmployees
  })

  function invalidateEmployees() {
    queryClient.invalidateQueries({ queryKey })
  }

  // actions
  const createMutation = useMutation({
    mutationFn: (req: CreateEmployeeRequest) => createEmployee(req),
    onSuccess: () => invalidateEmployees()
  })
  const updateMutation = useMutation({
    mutationFn: ({ id, req }: { id: number; req: UpdateEmployeeRequest }) =>
      updateEmployee(id, req),
    onSuccess: () => invalidateEmployees()
  })
  // 软删除
  const deleteMutation = useMutation({
    mutationFn: (id: number) => updateEmployee(id, { isActive: false }),
    onSuccess: () => invalidateEmployees()
  })
  return {
    employees: employeeQuery.data,
    isLoading: employeeQuery.isLoading,
    isError: employeeQuery.isError,

    createEmployee: createMutation.mutateAsync,
    updateEmployee: updateMutation.mutateAsync,
    deleteEmployee: deleteMutation.mutateAsync,

    isCreating: createMutation.isPending,
    isUpdating: updateMutation.isPending,
    isDeleting: deleteMutation.isPending
  }
}
