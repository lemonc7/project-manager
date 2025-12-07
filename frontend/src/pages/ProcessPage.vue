<template>
  <UDashboardPanel id="process">
    <template #header>
      <UDashboardNavbar title="项目流程">
        <template #leading>
          <UDashboardSidebarCollapse title="折叠侧栏" />
        </template>
        <template #right>
          <CreateProcessModal />
          <RefreshButton />
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <UTable
        sticky
        class="shrink-0"
        :data="processes ?? []"
        :columns="columns"
        :loading="isProcessLoading"
        :ui="{
          base: 'table-fixed border-separate border-spacing-0',
          thead: '[&>tr]:bg-elevated/50 [&>tr]:after:content-none',
          tbody: '[&>tr]:last:[&>td]:border-b-0',
          th: 'py-2 first:rounded-l-lg last:rounded-r-lg border-y border-default first:border-l last:border-r',
          td: 'py-2.5 border-b border-default',
          separator: 'h-0'
        }"
      />
    </template>
  </UDashboardPanel>
  <!-- 删除弹窗 -->
  <DeleteModal
    v-model:open="isDeleteModalOpen"
    :on-submit="handleDelete"
    title="删除流程"
    :description="`确认要删除[${selectedProcess?.name}]流程么？该操作无法撤销!`"
  />
</template>

<script setup lang="ts">
import { h, ref, resolveComponent } from "vue"
import { useProcesses } from "../composables/useProcesses"
import { type Row } from "@tanstack/table-core"
import type { Process, ProcessStatus } from "../models/processes"
import type { TableColumn } from "@nuxt/ui"
import useDeleteProcess from "../composables/useDeleteProcess"
import useUpdateProcess from "../composables/useUpdateProcess"
import EditableCell from "../components/EditableCell.vue"
import useEmployeeMap from "../composables/useEmployeeMap"

const toast = useToast()
const { data: processes, isLoading: isProcessLoading } = useProcesses()
const { employeeItems } = useEmployeeMap()

const isEditModalOpen = ref(false)
const isDeleteModalOpen = ref(false)
const selectedProcess = ref<Process | null>(null)

const UDropdownMenu = resolveComponent("UDropdownMenu")
const UBadge = resolveComponent("UBadge")
const USelectMenu = resolveComponent("USelectMenu")
const UButton = resolveComponent("UButton")

// 删除方法
const deleteProcess = useDeleteProcess({
  onSuccess: () => {
    toast.add({
      title: `流程 ${selectedProcess.value?.name} 删除成功`,
      color: "success"
    })
  },
  onError: () => {
    toast.add({
      title: "删除失败",
      color: "error"
    })
  }
})

function handleDelete() {
  if (selectedProcess.value) {
    deleteProcess.mutate(selectedProcess.value.id)
    isDeleteModalOpen.value = false
  }
}

const updateProcess = useUpdateProcess({
  onError: () => {
    toast.add({
      title: "更新失败",
      color: "error"
    })
  }
})

// 操作菜单定义
function getRowItems(row: Row<Process>) {
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
      onSelect: () => {
        isEditModalOpen.value = true
        selectedProcess.value = row.original
      }
    },
    {
      label: "删除",
      icon: "i-lucide-trash",
      color: "error",
      onSelect: () => {
        isDeleteModalOpen.value = true
        selectedProcess.value = row.original
      }
    }
  ]
}

const columns: TableColumn<Process>[] = [
  {
    accessorKey: "sequence",
    header: "序号"
  },
  {
    header: "名称",
    cell: ({ row }) => {
      return h(EditableCell, {
        value: row.original.name,
        onSave: (value: string) => {
          updateProcess.mutate({
            id: row.original.id,
            req: { name: value }
          })
        }
      })
    }
  },
  {
    header: "细节描述",
    cell: ({ row }) => {
      return h(EditableCell, {
        value: row.original.description,
        onSave: (value: string) => {
          updateProcess.mutate({
            id: row.original.id,
            req: { description: value }
          })
        }
      })
    }
  },
  {
    header: "负责人",
    cell: ({ row }) => {
      return h(USelectMenu, {
        modelValue: row.original.employeeID,
        "onUpdate:modelValue": (value: number) => {
          updateProcess.mutate({
            id: row.original.id,
            req: { employeeID: value }
          })
        },
        items: employeeItems.value,
        valueKey: "value",
        class: "w-28"
      })
    }
  },
  {
    header: "完成时间",
    cell: ({ row }) => {
      const process = row.original
      if (!process.finishedAt) return "-"

      const date = new Date(process.finishedAt)
      // 格式化为 YYYY-MM-DD HH:mm
      const time = date.toLocaleString("zh-CN", {
        year: "numeric",
        month: "2-digit",
        day: "2-digit",
        hour: "2-digit",
        minute: "2-digit",
        hour12: false
      })

      return process.spendTime ? `${time} (机时${process.spendTime}h)` : time
    }
  },
  {
    header: "状态",
    cell: ({ row }) => {
      const colorMap: Record<ProcessStatus, string> = {
        等待: "neutral",
        准备: "secondary",
        进行: "info",
        完成: "success",
        暂停: "warning",
        终止: "error"
      }
      const color = colorMap[row.original.status]
      return h(
        UBadge,
        { class: "capitalize", color, variant: "subtle" },
        () => row.original.status
      )
    }
  },
  {
    header: "备注",
    cell: ({ row }) => {
      return h(EditableCell, {
        value: row.original.remark,
        onSave: (value: string) => {
          updateProcess.mutate({
            id: row.original.id,
            req: { remark: value }
          })
        }
      })
    }
  },
  {
    id: "actions",
    cell: ({ row }) => {
      return h(
        "div",
        {
          class: "text-right"
        },
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
</script>
