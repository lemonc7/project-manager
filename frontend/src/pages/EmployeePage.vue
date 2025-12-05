<template>
  <!-- actions下拉菜单 -->
  <DefineTemplate v-slot="{ employee, onEdit, onDelete }">
    <UDropdownMenu
      :content="{ align: 'end' }"
      :items="[
        {
          label: '编辑',
          icon: 'i-lucide-edit',
          onSelect: () => onEdit(employee)
        },
        {
          label: '删除',
          icon: 'i-lucide-trash',
          color: 'error',
          onSelect: () => onDelete(employee)
        }
      ]"
    >
      <UButton variant="outline" icon="i-lucide-menu" />
    </UDropdownMenu>
  </DefineTemplate>
  <!-- table渲染 -->
  <UTable
    sticky
    :loading="isEmployeesLoading || isTeamsLoading"
    :data="employees"
    :columns="columns"
    class="flex-1"
  />
</template>

<script setup lang="ts">
import type { TableColumn } from "@nuxt/ui"
import useEmployees from "../composables/useEmployees"
import { RoleMap, type Employee } from "../models/employees"
import getTeamName from "../utils/getTeamName"
import useTeams from "../composables/useTeams"
import { createReusableTemplate } from "@vueuse/core"
import { h, ref } from "vue"

const [DefineTemplate, ReuseTemplate] = createReusableTemplate<{
  employee: Employee
  onEdit: (employee: Employee) => void
  onDelete: (employee: Employee) => void
}>()

const { data: employees, isLoading: isEmployeesLoading } = useEmployees()
const { data: teams, isLoading: isTeamsLoading } = useTeams()

const isEditModalOpen = ref(false)
const isDeleteModalOpen = ref(false)
const isCreateModalOpen = ref(false)
const selectedEmployee = ref<Employee | null>(null)

const columns: TableColumn<Employee>[] = [
  {
    accessorKey: "id",
    header: "ID"
  },
  {
    accessorKey: "name",
    header: "成员"
  },
  {
    header: "小组",
    cell: ({ row }) => {
      const team = getTeamName(row.original.teamID, teams.value)
      return team
    }
  },
  {
    header: "角色",
    cell: ({ row }) => {
      const employee = row.original
      const labels = employee.roles?.map((r) => RoleMap[r] ?? r).join("，")
      return labels
    }
  },
  {
    id: "actions",
    header: "操作",
    cell: ({ row }) => {
      const employee = row.original
      return h(ReuseTemplate, {
        employee: employee,
        onEdit: (employee: Employee) => {
          isEditModalOpen.value = true
          selectedEmployee.value = employee
        },
        onDelete: (employee: Employee) => {
          isDeleteModalOpen.value = true
          selectedEmployee.value = employee
        }
      })
    }
  }
]
</script>
