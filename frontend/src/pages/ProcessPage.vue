<template>
  <DefineTemplate v-slot="{ process }">
    <USelectMenu
      :model-value="process.employeeID"
      :items="employeeItems"
      @update:model-value="(value:number) => handleUpdate({id:process.id,req:{employeeID:value}}) "
      value-key="value"
      placeholder="请选择"
      class="min-w-24"
    />
  </DefineTemplate>
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
      <div class="flex flex-wrap items-center justify-between gap-2">
        <UFieldGroup>
          <UButton :label="gotProject?.name" variant="outline" />
          <UButton :label="gotProject?.description" variant="outline" />
          <UButton
            :label="
              gotProject?.remainingDays >= 0
                ? `剩余${gotProject?.remainingDays}天`
                : `超时${Math.abs(gotProject?.remainingDays)}天`
            "
            variant="outline"
            :color="gotProject.remainingDays <= 0 ? 'error' : 'info'"
            v-if="gotProject?.remainingDays !== undefined"
          />
        </UFieldGroup>

        <USelectMenu
          v-model="employeeFilter"
          icon="i-lucide-user"
          value-key="value"
          :items="[{ label: '全部', value: 0 }, ...employeeItems]"
          :ui="{
            trailingIcon:
              'group-date-[state=open]:rotate-180 transition-transform duration-200'
          }"
          class="min-w-30"
          :search-input="{
            placeholder: '过滤负责人...'
          }"
        />
      </div>
      <UTable
        sticky="header"
        ref="table"
        v-model:row-selection="rowSelection"
        :data="processes"
        :columns="columns"
        :loading="isLoading"
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
  <AlertModal
    v-model:open="isDeleteModalOpen"
    :on-submit="handleDelete"
    title="删除流程"
    label="删除"
    color="error"
    :description="`确认要删除[${selectedProcess?.name}]流程么？该操作无法撤销!`"
  />
  <!-- 编辑弹窗 -->
  <EditProcessModal v-model:open="isEditModalOpen" :process="selectedProcess" />
  <!-- 完成弹窗 -->
  <ReworkOrFinishProcessDrawer
    v-model:open="isFinishDrawerOpen"
    :process="selectedProcess"
    type="finish"
  />
  <!-- 返工弹窗 -->
  <ReworkOrFinishProcessDrawer
    v-model:open="isReworkDrawerOpen"
    :process="selectedProcess"
    type="rework"
  />
</template>

<script setup lang="ts">
import { useTemplateRef, h, nextTick, ref, resolveComponent, watch } from "vue"
import { type Row, type Table } from "@tanstack/table-core"
import type {
  Process,
  ProcessStatus,
  UpdateProcessRequest
} from "../models/processes"
import type { TableColumn } from "@nuxt/ui"
import EditableCell from "../components/EditableCell.vue"
import { useEmployeeStore } from "../composables/useEmployeeStore"
import { storeToRefs } from "pinia"
import { useProcessStore } from "../composables/useProcessStore"
import { createReusableTemplate } from "@vueuse/core"
import { useProjectStore } from "../composables/useProjectStore"

const toast = useToast()
const table = useTemplateRef<{ tableApi: Table<Process> }>("table")

const employeeStore = useEmployeeStore()
const processStore = useProcessStore()
const projectStore = useProjectStore()
const { employeeItems } = storeToRefs(employeeStore)
const { processes, isLoading, sequenceToScroll } = storeToRefs(processStore)
const { deleteProcess, updateProcess, clearSequenceToScroll } = processStore
const { gotProject } = storeToRefs(projectStore)

const isEditModalOpen = ref(false)
const isDeleteModalOpen = ref(false)
const isFinishDrawerOpen = ref(false)
const isReworkDrawerOpen = ref(false)
const selectedProcess = ref<Process | null>(null)

const UDropdownMenu = resolveComponent("UDropdownMenu")
const UBadge = resolveComponent("UBadge")
const USelectMenu = resolveComponent("USelectMenu")
const UButton = resolveComponent("UButton")
const UCheckbox = resolveComponent("UCheckbox")

// 负责人选择框复用
const [DefineTemplate, ReuseTemplate] = createReusableTemplate<{
  process: Process
}>()

const rowSelection = ref<Record<number, boolean>>({})

// 删除方法
async function handleDelete() {
  if (selectedProcess.value) {
    try {
      await deleteProcess(selectedProcess.value.id)
      toast.add({
        title: `流程 ${selectedProcess.value.name} 删除成功`,
        color: "success"
      })
    } catch (error) {
      toast.add({
        title: "删除失败",
        color: "error"
      })
      console.log("删除流程失败:", error)
    } finally {
      isDeleteModalOpen.value = false
    }
  }
}

// 更新
async function handleUpdate({
  id,
  req
}: {
  id: number
  req: UpdateProcessRequest
}) {
  try {
    await updateProcess({ id, req })
  } catch (err) {
    toast.add({
      title: "更新失败",
      color: "error"
    })
    console.log("更新流程失败:", err)
  }
}

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
      label: "完成",
      icon: "i-lucide-check",
      disabled: row.original.status !== "进行",
      onSelect: () => {
        isFinishDrawerOpen.value = true
        selectedProcess.value = row.original
      }
    },
    {
      label: "返工",
      icon: "i-lucide-redo",
      disabled: row.original.status !== "进行",
      onSelect: () => {
        isReworkDrawerOpen.value = true
        selectedProcess.value = row.original
      }
    },
    {
      type: "separator"
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
    id: "select",
    header: ({ table }) =>
      h(UCheckbox, {
        modelValue: table.getIsSomePageRowsSelected()
          ? "indeterminate"
          : table.getIsAllPageRowsSelected(),
        "onUpdate:modelValue": (value: boolean | "indeterminate") =>
          table.toggleAllPageRowsSelected(!!value),
        ariaLabel: "Select all"
      }),
    cell: ({ row }) =>
      h(UCheckbox, {
        modelValue: row.getIsSelected(),
        "onUpdate:modelValue": (value: boolean | "indeterminate") =>
          row.toggleSelected(!!value),
        ariaLabel: "Select row"
      })
  },
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
          handleUpdate({
            id: row.original.id,
            req: { name: value }
          })
        }
      })
    }
  },
  {
    header: "描述",
    cell: ({ row }) => {
      return h(EditableCell, {
        value: row.original.description,
        onSave: (value: string) => {
          handleUpdate({
            id: row.original.id,
            req: { description: value }
          })
        }
      })
    }
  },
  {
    id: "employee",
    header: "负责人",
    accessorFn: (row) => row.employeeID,
    filterFn: "equals",
    cell: ({ row }) => {
      return h(ReuseTemplate, {
        process: row.original
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
        返工: "warning",
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
          handleUpdate({
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

// 负责人过滤
const employeeFilter = ref(0)
watch(employeeFilter, (newVal) => {
  if (!table.value?.tableApi) return
  const employeeColumn = table.value.tableApi.getColumn("employee")
  if (!employeeColumn) return

  if (newVal === 0) {
    employeeColumn.setFilterValue(undefined)
  } else {
    employeeColumn.setFilterValue(newVal)
  }
})

// 滚动到选择的行
watch(
  [processes, sequenceToScroll],
  async ([newVal, targetSequence]) => {
    if (targetSequence !== null && newVal && newVal.length > 0) {
      const selectedIndex = targetSequence - 1
      if (selectedIndex >= 0 && selectedIndex < newVal.length) {
        // 设置选中状态
        rowSelection.value = { [selectedIndex]: true }

        // 等table渲染完成
        await nextTick()

        // 查找目标行
        const rows = document.querySelectorAll("tbody tr")
        const targetRow = rows[selectedIndex] as HTMLElement

        // 滚动
        if (targetRow) {
          targetRow.scrollIntoView({
            behavior: "smooth",
            block: "nearest"
          })
        }

        // 滚动完成，清理store中的sequence
        clearSequenceToScroll()
      }
    }
  },
  { immediate: true }
)
</script>
