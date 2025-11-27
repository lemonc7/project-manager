<template>
  <div></div>
</template>

<script setup lang="ts">
import EditableCell from "@/components/common/EditableCell.vue"
import DropdownMenu from "@/components/ui/dropdown-menu/DropdownMenu.vue"
import useDeleteEmployee from "@/composables/use-delete-employee"
import useEmployees from "@/composables/use-employees"
import useTeams from "@/composables/use-teams"
import useUpdateEmployee from "@/composables/use-update-employee"
import type { Employee } from "@/model/employees"
import type { ColumnDef } from "@tanstack/vue-table"
import { h, ref } from "vue"

const { data: employees, isLoading, isError } = useEmployees()
const { data: teams } = useTeams()

const isDialogOpen = ref(false)
const selectedEmployee = ref<Employee | null>(null)

const deleteEmployee = useDeleteEmployee({
  onSuccess: () => {
    isDialogOpen.value = false
    selectedEmployee.value = null
  },
  onError: () => {
    isDialogOpen.value = false
    selectedEmployee.value = null
  }
})

const updateEmployee = useUpdateEmployee({})

function handleDeleteConfirm() {
  if (selectedEmployee.value) {
    deleteEmployee.mutate(selectedEmployee.value.id)
  }
}

// 表格定义
const columns: ColumnDef<Employee>[] = [
  {
    accessorKey: "name",
    header: "成员"
  },
  {
    header: "小组",
    cell: ({ row }) => {
      return h(EditableCell, {
        value: row.original.teamID,
        type: "select",
        options: teams.value?.map((t) => ({
          value: t.id,
          label: t.name
        })),
        onSave: (value: number) =>
          updateEmployee.mutate({
            id: row.original.id,
            req: { teamID: value }
          })
      })
    }
  },
  {
    id: "actions",
    header: "操作",
    size: 30,
    cell: ({ row }) => {

    }
  }
]
</script>
