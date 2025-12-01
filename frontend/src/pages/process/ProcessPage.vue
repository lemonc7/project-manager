<template>
  <div
    v-if="isProcessesError || isEmployeesError"
    class="flex items-center justify-center h-full"
  >
    <span>加载失败</span>
  </div>
  <div
    v-else-if="isProcessesLoading || isEmployeesLoading"
    class="flex items-center justify-center h-full"
  >
    <span>流程加载中...</span>
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
              <TableRow :date-state="row.getIsSelected() && 'selected'">
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
  <!-- Header中的按钮 -->
  <SafeTeleport to="#header-actions">
    <Button @click="isCreateDialogOpen = true">创建流程</Button>
  </SafeTeleport>
  <!-- 创建弹窗 -->
  <CreateProcessDialog
    v-model:open="isCreateDialogOpen"
    :max-sequence="maxSequence"
    :employees="employees ?? []"
  />
  <!-- 编辑弹窗 -->
  <EditProcessDialog
    v-model:open="isUpdateDialogOpen"
    :employees="employees ?? []"
    :process="selectedProcess"
  />
  <!-- 删除弹窗 -->
  <AlertDialog v-model:open="isDeleteDialogOpen">
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle> 删除流程 </AlertDialogTitle>
        <AlertDialogDescription>
          确定要删除
          <span class="font-bold text-red-500">{{
            selectedProcess?.name
          }}</span>
          工艺吗? 这个操作不能撤销，将从数据库中删除该工艺流程。
        </AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogCancel
          @click="selectedProcess = null"
          :disabled="deleteProcess.isPending.value"
        >
          取消
        </AlertDialogCancel>
        <AlertDialogAction
          @click="handleDeleteConfirm"
          :disabled="deleteProcess.isPending.value"
        >
          {{ deleteProcess.isPending.value ? "删除中..." : "确认" }}
        </AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
  <!-- 完成弹窗 -->
  <Drawer v-model:open="isFinishDialogOpen">
    <DrawerContent>
      <div class="mx-auto w-full max-w-sm">
        <DrawerHeader>
          <DrawerTitle>流程设置</DrawerTitle>
          <DrawerDescription
            >将选中的流程设置为完成状态，并设置机时</DrawerDescription
          >
        </DrawerHeader>
        <div class="p-4 pb-0">
          <div class="flex items-center justify-center space-x-2">
            <Button
              variant="outline"
              size="icon"
              class="h-8 w-8 shrink-0 rounded-full"
              :disabled="spendTime <= 0"
              @click="spendTime -= 0.5"
            >
              <Minus />
            </Button>
            <div class="flex-1 text-center">
              <div class="text-7xl font-bold tracking-tighter">
                <!-- 展示模式 -->
                <div
                  v-if="!editing"
                  class="cursor-pointer"
                  @click="startEdit"
                >
                  {{ spendTime }}
                </div>
                <!-- 输入模式 -->
                <input
                  v-else
                  v-model.number="spendTime"
                  ref="inputRef"
                  @blur="finishEdit"
                  @keyup.enter="finishEdit"
                  type="number"
                  step="0.5"
                  min="0"
                  class="w-32 text-center bg-transparent outline-none text-5xl font-bold"
                />
              </div>
              <div class="text-muted-foreground text-sm">小时</div>
            </div>
            <Button
              variant="outline"
              size="icon"
              class="h-8 w-8 shrink-0 rounded-full"
              @click="spendTime += 0.5"
            >
              <Plus />
            </Button>
          </div>
        </div>
        <DrawerFooter>
          <Button @click="handleFinishConfirm">确认</Button>
          <DrawerClose as-child>
            <Button variant="outline">取消</Button>
          </DrawerClose>
        </DrawerFooter>
      </div>
    </DrawerContent>
  </Drawer>
</template>

<script setup lang="ts">
import EditableCell from "@/components/common/EditableCell.vue"
import useEmployees from "@/composables/use-employees"
import { useProcesses } from "@/composables/use-processes"
import useUpdateProcess from "@/composables/use-update-process"
import type { Process } from "@/models/processes"
import { hasRole } from "@/util/auth"
import {
  FlexRender,
  getCoreRowModel,
  useVueTable,
  type ColumnDef
} from "@tanstack/vue-table"
import { computed, h, nextTick, ref, watch } from "vue"
import { toast } from "vue-sonner"
import ProcessStatus from "./ProcessStatus.vue"
import ProcessActions from "./ProcessActions.vue"
import useDeleteProcess from "@/composables/use-delete-process"
import useFinishProcess from "@/composables/use-finish-process"
import Button from "@/components/ui/button/Button.vue"
import { Minus, Plus } from "lucide-vue-next"

const {
  data: processes,
  isLoading: isProcessesLoading,
  isError: isProcessesError
} = useProcesses()

const {
  data: employees,
  isLoading: isEmployeesLoading,
  isError: isEmployeesError
} = useEmployees()

const isCreateDialogOpen = ref(false)
const isUpdateDialogOpen = ref(false)
const isDeleteDialogOpen = ref(false)
const isFinishDialogOpen = ref(false)

const selectedProcess = ref<Process | null>(null)

const updateProcess = useUpdateProcess({
  onError: () => {
    toast.error("更新失败")
  }
})

// 编辑
const handleEditProcess = (process: Process) => {
  selectedProcess.value = process
  isUpdateDialogOpen.value = true
}

// 删除
const deleteProcess = useDeleteProcess({
  onSuccess: () => {
    toast.success(`流程 ${selectedProcess.value?.name} 删除成功`)
    isDeleteDialogOpen.value = false
    selectedProcess.value = null
  },
  onError: () => {
    isDeleteDialogOpen.value = false
    selectedProcess.value = null
    toast.error("删除失败")
  }
})

const handleDeleteConfirm = (e: MouseEvent) => {
  e.preventDefault()
  if (selectedProcess.value) {
    deleteProcess.mutate(selectedProcess.value.id)
  }
}

// 完成
const spendTime = ref(0)
const editing = ref(false)
const inputRef = ref<HTMLInputElement>()

const startEdit = () => {
  editing.value = true
  nextTick(() => {
    inputRef.value?.focus()
    inputRef.value?.select()
  })
}

const finishEdit = () => {
  if (!spendTime.value && spendTime.value !== 0) {
    spendTime.value = 0
  }
  editing.value = false
}

const finishProcess = useFinishProcess({
  onSuccess: () => {
    toast.success("操作成功")
    isFinishDialogOpen.value = false
    selectedProcess.value = null
  },
  onError: () => {
    isFinishDialogOpen.value = false
    selectedProcess.value = null
    toast.error("操作失败")
  }
})

const handleFinishConfirm = (e: MouseEvent) => {
  e.preventDefault()
  if (selectedProcess.value) {
    finishProcess.mutate({
      id: selectedProcess.value.id,
      req: { spendTime: spendTime.value }
    })
  }
}

// 监听isFinishDialogOpen变化，关闭时重置spendTime
watch(isFinishDialogOpen, (newVal) => {
  if (!newVal) {
    spendTime.value = 0
  }
})

// 计算processes里的最大序号
const maxSequence = computed(() => {
  if (processes.value?.length) {
    const maxVal = processes.value.reduce((max, p) => {
      return Math.max(max, p.sequence)
    }, 0)
    return maxVal + 1
  }
  return 1
})

const columns: ColumnDef<Process>[] = [
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
    header: "描述",
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
      return h(EditableCell, {
        value: row.original.employeeID,
        type: "select",
        options: employees.value?.map((e) => ({
          value: e.id,
          label: e.name
        })),
        onSave: (value: number) => {
          updateProcess.mutate({
            id: row.original.id,
            req: { employeeID: value }
          })
        },
        disabled: !hasRole("employee", "admin")
      })
    }
  },
  {
    header: "完成时间",
    cell: ({ row }) => {
      const process = row.original
      if (!process.finishedAt) {
        return "-"
      }
      let time = new Date(process.finishedAt).toISOString().split("T")[0]
      if (process.spendTime) {
        time += ` (机时${process.spendTime}h)`
      }
      return time
    }
  },
  {
    header: "状态",
    cell: ({ row }) => {
      return h(ProcessStatus, {
        status: row.original.status,
        onClick: () => {
          selectedProcess.value = row.original
          isFinishDialogOpen.value = true
        }
      })
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
    header: "操作",
    size: 30,
    cell: ({ row }) => {
      const process = row.original
      return h(ProcessActions, {
        process: process,
        onEdit: handleEditProcess,
        onDelete: (p) => {
          selectedProcess.value = p
          isDeleteDialogOpen.value = true
        }
      })
    }
  }
]

const table = useVueTable({
  get data() {
    return processes.value ?? []
  },
  columns,
  getCoreRowModel: getCoreRowModel()
})
</script>
