<template>
  <UDashboardPanel id="employee">
    <template #header>
      <UDashboardNavbar title="成员列表">
        <template #leading>
          <UDashboardSidebarCollapse title="折叠侧栏" />
        </template>
        <template #right>
          <div class="flex flex-wrap items-center gap-2">
            <CreateEmployeeModal />
            <RefreshButton />
          </div>
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <div class="flex flex-wrap items-center justify-between gap-2">
        <UInput
          v-model="name"
          class="max-w-sm"
          icon="i-lucide-search"
          placeholder="姓名过滤..."
        />
        <USelect
          v-model="teamFilter"
          icon="material-symbols:family-group"
          :items="[{ label: '全部', value: 0 }, ...teamItems]"
          :ui="{
            trailingIcon:
              'group-data-[state=open]:rotate-180 transition-transform duration-200'
          }"
          placeholder="小组过滤..."
          class="min-w-40"
        />
      </div>

      <!-- table渲染 -->
      <UTable
        sticky="header"
        ref="table"
        v-model:column-filters="columnFilters"
        :loading="isLoading"
        :data="employees"
        :columns="columns"
        :ui="{
          base: 'table-fixed border-separate border-spacing-0',
          thead: '[&>tr]:bg-elevated/50 [&>tr]:after:content-none',
          tbody: '[&>tr]:last:[&>td]:border-b-0',
          th: 'py-2 first:rounded-l-lg last:rounded-r-lg border-y border-default first:border-l last:border-r',
          td: 'py-2 border-b border-default',
          separator: 'h-0'
        }"
      />
    </template>
  </UDashboardPanel>
  <!-- 删除弹窗 -->
  <AlertModal
    v-model:open="isDeleteModalOpen"
    :on-submit="handleDelete"
    title="删除成员"
    label="删除"
    color="error"
    :description="`确认要删除[${selectedEmployee?.name}]么？该操作无法撤销!`"
  />
  <!-- 编辑弹窗 -->
  <EditEmployeeModal
    v-model:open="isEditModalOpen"
    :employee="selectedEmployee"
  />
  <!-- 重置密码弹窗 -->
  <ResetPasswordModal
    v-model:open="isResetPasswordModalOpen"
    :employee="selectedEmployee"
  />
</template>

<script setup lang="ts">
import type { TableColumn } from "@nuxt/ui"
import { type Employee } from "../models/employees"
import { h, ref, resolveComponent, useTemplateRef, watch } from "vue"
import { type Row, type Table } from "@tanstack/table-core"
import getRoleLabels from "../utils/getRoleLabels"
import { useEmployeeStore } from "../composables/useEmployeeStore"
import { useTeamStore } from "../composables/useTeamStore"
import { storeToRefs } from "pinia"
import { hasAtLeastRole } from "../utils/auth"

// 数据获取
const teamStore = useTeamStore()
const employeeStore = useEmployeeStore()
const { deleteEmployee } = employeeStore
const { employees, isLoading } = storeToRefs(employeeStore)
const { getTeamName } = teamStore
const { teamItems } = storeToRefs(teamStore)

const toast = useToast()

// table定义
const table = useTemplateRef<{ tableApi: Table<Employee> }>("table")

// 模态框
const isEditModalOpen = ref(false)
const isDeleteModalOpen = ref(false)
const isResetPasswordModalOpen = ref(false)
const selectedEmployee = ref<Employee | null>(null)

async function handleDelete() {
  if (selectedEmployee.value) {
    try {
      await deleteEmployee(selectedEmployee.value.id)
      toast.add({
        title: `成员 ${selectedEmployee.value.name} 删除成功`,
        color: "success"
      })
    } catch (error) {
      toast.add({
        title: "删除失败",
        color: "error"
      })
      console.log("删除成员失败:", error)
    } finally {
      isDeleteModalOpen.value = false
    }
  }
}

// 组件定义
const UDropdownMenu = resolveComponent("UDropdownMenu")
const UButton = resolveComponent("UButton")

// 操作菜单定义
function getRowItems(row: Row<Employee>) {
  return [
    {
      type: "label",
      label: "操作"
    },
    {
      type: "separator"
    },
    {
      label: "编辑",
      icon: "i-lucide-pencil",
      disabled: !hasAtLeastRole("admin"),
      onSelect: () => {
        isEditModalOpen.value = true
        selectedEmployee.value = row.original
      }
    },
    {
      label: "重置密码",
      icon: "i-lucide-key",
      disabled: !hasAtLeastRole("admin"),
      onSelect: () => {
        isResetPasswordModalOpen.value = true
        selectedEmployee.value = row.original
      }
    },
    {
      type: "separator"
    },
    {
      label: "删除",
      icon: "i-lucide-trash",
      color: "error",
      disabled: !hasAtLeastRole("admin"),
      onSelect: () => {
        isDeleteModalOpen.value = true
        selectedEmployee.value = row.original
      }
    }
  ]
}

// table列定义
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
    id: "team",
    header: "小组",
    accessorFn: (row) => row.teamID,
    filterFn: "equals",
    cell: ({ row }) => getTeamName(row.original.teamID)
  },
  {
    header: "角色",
    cell: ({ row }) => getRoleLabels(row.original.roles)
  },
  {
    id: "actions",
    cell: ({ row }) => {
      return h(
        "div",
        { class: "text-right" },
        h(
          UDropdownMenu,
          {
            content: {
              align: "end"
            },
            items: getRowItems(row)
          },
          () =>
            h(UButton, {
              icon: "i-lucide-ellipsis-vertical",
              variant: "ghost",
              class: "ml-auto"
            })
        )
      )
    }
  }
]

// table姓名过滤
const columnFilters = ref([
  {
    id: "name",
    value: ""
  }
])

// table小组过滤
const teamFilter = ref(0)

// 监视小组选择
watch(
  () => teamFilter.value,
  (newVal) => {
    if (!table.value?.tableApi) return
    const statusColumn = table.value.tableApi.getColumn("team")
    if (!statusColumn) return

    if (newVal === 0) {
      statusColumn.setFilterValue(undefined)
    } else {
      statusColumn.setFilterValue(newVal)
    }
  }
)

// 监视姓名过滤
const name = ref("")
let timer: ReturnType<typeof setTimeout> | null = null
watch(
  name,
  (newVal) => {
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => {
      table.value?.tableApi.getColumn("name")?.setFilterValue(newVal)
      timer = null
    }, 300)
  },
  { flush: "post" }
)
</script>
