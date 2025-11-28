<template>
  <div
    v-if="isError || isTeamsError"
    class="flex items-center justify-center h-full"
  >
    <span>加载失败</span>
  </div>
  <div v-else-if="isLoading || isTeamsLoading">
    <span>成员列表加载中...</span>
  </div>

  <!-- 表格渲染 -->
  <div v-else>
    <div class="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow
            v-for="headerGroup in table.getHeaderGroups()"
            :key="headerGroup.id"
          >
            <TableHead v-for="header in headerGroup.headers" :key="header.id">
              <FlexRender
                v-if="!header.isPlaceholder"
                :render="header.column.columnDef.header"
                :props="header.getContext()"
              />
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <template v-if="table.getRowModel().rows?.length">
            <template v-for="row in table.getRowModel().rows" :key="row.id">
              <TableRow :data-state="row.getIsSelected() && 'selected'">
                <TableCell v-for="cell in row.getVisibleCells()" :key="cell.id">
                  <FlexRender
                    :render="cell.column.columnDef.cell"
                    :props="cell.getContext()"
                  />
                </TableCell>
              </TableRow>
            </template>
          </template>
          <TableRow v-else>
            <TableCell :colspan="columns.length" class="h-24 text-center">
              无数据
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>
  </div>

  <!-- 弹窗 -->
  <AlertDialog v-model:open="isDialogOpen">
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle> 删除成员 </AlertDialogTitle>
        <AlertDialogDescription>
          确定要删除{{ selectedEmployee?.name }}吗?
          这个操作不能撤销，将从数据库中删除该成员。
        </AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogCancel
          @click="selectedEmployee = null"
          :disabled="deleteEmployee.isPending.value"
        >
          取消
        </AlertDialogCancel>
        <AlertDialogAction
          @click="handleDeleteConfirm"
          :disabled="deleteEmployee.isPending.value"
        >
          {{ deleteEmployee.isPending.value ? "删除中..." : "确认" }}
        </AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>

<script setup lang="ts">
import EditableCell from "@/components/common/EditableCell.vue"
import useDeleteEmployee from "@/composables/use-delete-employee"
import useEmployees from "@/composables/use-employees"
import useTeams from "@/composables/use-teams"
import useUpdateEmployee from "@/composables/use-update-employee"
import type { Employee } from "@/model/employees"
import {
  FlexRender,
  getCoreRowModel,
  useVueTable,
  type ColumnDef
} from "@tanstack/vue-table"
import { h, ref } from "vue"
import EmployeeActions from "./EmployeeActions.vue"
import AlertDialog from "@/components/ui/alert-dialog/AlertDialog.vue"
import AlertDialogContent from "@/components/ui/alert-dialog/AlertDialogContent.vue"
import AlertDialogHeader from "@/components/ui/alert-dialog/AlertDialogHeader.vue"
import AlertDialogTitle from "@/components/ui/alert-dialog/AlertDialogTitle.vue"
import AlertDialogDescription from "@/components/ui/alert-dialog/AlertDialogDescription.vue"
import AlertDialogFooter from "@/components/ui/alert-dialog/AlertDialogFooter.vue"
import AlertDialogCancel from "@/components/ui/alert-dialog/AlertDialogCancel.vue"
import AlertDialogAction from "@/components/ui/alert-dialog/AlertDialogAction.vue"

const { data: employees, isLoading, isError } = useEmployees()
const {
  data: teams,
  isLoading: isTeamsLoading,
  isError: isTeamsError
} = useTeams()

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

function handleDeleteConfirm(e: MouseEvent) {
  e.preventDefault()
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
      const employee = row.original
      return h(EmployeeActions, {
        employee,
        onEdit: (emp) => console.log("编辑", emp.id),
        onDelete: (emp) => {
          selectedEmployee.value = emp
          isDialogOpen.value = true
        }
      })
    }
  }
]

const table = useVueTable({
  get data() {
    return employees.value ?? []
  },
  columns,
  getCoreRowModel: getCoreRowModel()
})
</script>
