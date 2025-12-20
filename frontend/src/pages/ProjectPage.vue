<template>
  <UDashboardPanel id="project">
    <template #header>
      <UDashboardNavbar title="项目列表">
        <template #leading>
          <UDashboardSidebarCollapse title="折叠侧栏" />
        </template>
        <template #right>
          <div class="flex flex-wrap items-center gap-2">
            <CreateProjectModal />
            <RefreshButton />
          </div>
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <div class="flex flex-wrap items-center justify-between gap-2">
        <!-- 项目名称过滤 -->
        <UInput
          v-model="name"
          :class="isMobile ? 'w-36' : 'max-w-sm'"
          icon="i-lucide-search"
          placeholder="过滤项目名称..."
        />
        <div class="flex flex-wrap items-center gap-2">
          <UFieldGroup>
            <!-- 起始时间 -->
            <UTooltip text="起始日期">
              <UInputDate v-model="startDate" ref="inputStartDate">
                <template #trailing>
                  <UPopover :reference="inputStartDate?.inputsRef[3]?.$el">
                    <UButton
                      color="neutral"
                      variant="link"
                      size="sm"
                      icon="i-lucide-calendar"
                      aria-label="选择起始日期"
                      class="px-0"
                    />
                    <template #content>
                      <UCalendar v-model="startDate" class="p-2" />
                    </template>
                  </UPopover>
                </template>
              </UInputDate>
            </UTooltip>
            <span
              class="flex items-center select-none px-1 text-black relative -top-px"
            >
              —
            </span>
            <!-- 终止时间 -->
            <UTooltip text="终止日期">
              <UInputDate v-model="endDate" ref="inputEndDate">
                <template #trailing>
                  <UPopover :reference="inputEndDate?.inputsRef[3]?.$el">
                    <UButton
                      color="neutral"
                      variant="link"
                      size="sm"
                      icon="i-lucide-calendar"
                      aria-label="选择终止日期"
                      class="px-0"
                    />
                    <template #content>
                      <UCalendar v-model="endDate" class="p-2" />
                    </template>
                  </UPopover>
                </template>
              </UInputDate>
            </UTooltip>
          </UFieldGroup>
          <UFieldGroup>
            <!-- 状态过滤 -->
            <UTooltip text="状态">
              <USelect
                v-model="projectState.status"
                multiple
                placeholder="过滤状态..."
                icon="i-lucide-list-todo"
                :items="[...ProjectStatusOptions]"
                :ui="{
                  trailingIcon:
                    'group-date-[state=open]:rotate-180 transition-transform duration-200'
                }"
                :class="isMobile ? 'w-28' : 'min-w-44'"
              />
            </UTooltip>
            <!-- 负责人过滤 -->
            <UTooltip text="负责人">
              <USelectMenu
                v-model="projectState.managerID"
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
            </UTooltip>
            <UTooltip text="重置过滤条件">
              <UButton
                icon="i-lucide-list-restart"
                variant="outline"
                @click="handleResetFilter"
              />
            </UTooltip>
            <UTooltip text="过滤项目">
              <UButton
                icon="i-lucide-check"
                color="neutral"
                @click="handleSetFilter"
              />
            </UTooltip>
          </UFieldGroup>
        </div>
      </div>
      <!-- table渲染 -->
      <UTable
        id="project-table"
        ref="table"
        sticky="header"
        v-model:column-filters="columnFilters"
        v-model:row-selection="rowSelection"
        :column-visibility="{
          id: !isMobile,
          startedAt: !isMobile,
          remark: !isMobile
        }"
        :data="projects"
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
    title="删除项目"
    label="删除"
    color="error"
    :description="`确认要删除[${selectedProject?.name}]项目么？该操作还会删除相关流程，无法撤销！`"
  />
  <!-- 编辑弹窗 -->
  <EditProjectModal v-model:open="isEditModalOpen" :project="selectedProject" />
  <!-- 开始弹窗 -->
  <AlertModal
    v-model:open="isStartModalOpen"
    :on-submit="handleStart"
    title="开始项目"
    label="开始"
    color="neutral"
    description="将项目状态设置为进行中，并同步设置流程状态。"
  />
  <!-- 完成弹窗 -->
  <AlertModal
    v-model:open="isFinishModalOpen"
    :on-submit="handleFinish"
    title="完成项目"
    label="完成"
    color="neutral"
    description="请确认项目已完成，该项目的所有流程在任务面板会隐藏。"
  />
  <!-- 重置弹窗 -->
  <AlertModal
    v-model:open="isResetModalOpen"
    :on-submit="handleReset"
    title="重置项目"
    label="重置"
    color="error"
    :description="`确认要重置[${selectedProject?.name}]项目么？该操作会重置项目和相关流程的状态和时间，无法撤销！`"
  />
  <!-- 复制弹窗 -->
  <CopyProjectModal v-model:open="isCopyModalOpen" :project="selectedProject" />
</template>

<script setup lang="ts">
import {
  computed,
  h,
  ref,
  resolveComponent,
  useTemplateRef,
  watch
} from "vue"
import { type Row } from "@tanstack/table-core"
import {
  ProjectStatusOptions,
  type ListProjectsRequest,
  type Project,
  type ProjectStatus
} from "../models/projects"
import { useRouter } from "vue-router"
import type { TableColumn } from "@nuxt/ui"
import { storeToRefs } from "pinia"
import { useQueryClient } from "@tanstack/vue-query"
import { hasAtLeastRole } from "../utils/auth"
import { useEmployeeStore } from "../stores/employee"
import { useProjects } from "../composables/useProjects"
import { useScreens } from "../composables/useScreens"
import {
  getDefaultProjectState,
  getProjectState,
  resetProjectState,
  setProjectState
} from "../utils/projectState"
import { parseDate, CalendarDate } from "@internationalized/date"
import { toRFC3339 } from "../utils/timeConvert"
import { exportProject } from "../api/projects"

defineOptions({
  name: "project"
})

const router = useRouter()
const toast = useToast()
const queryClient = useQueryClient()

const {
  projects,
  isLoading,
  setQuery,
  deleteProject,
  startProject,
  finishProject,
  resetProject
} = useProjects()
const { isMobile } = useScreens()

const projectState = ref<ListProjectsRequest>(getProjectState())

const inputStartDate = useTemplateRef("inputStartDate")
const inputEndDate = useTemplateRef("inputEndDate")

const startDate = computed({
  get() {
    try {
      return projectState.value.startAt
        ? parseDate(projectState.value.startAt)
        : null
    } catch {
      return null
    }
  },
  set(val: CalendarDate | null) {
    projectState.value.startAt = val ? val.toString() : undefined
  }
})

const endDate = computed({
  get() {
    try {
      return projectState.value.endAt
        ? parseDate(projectState.value.endAt)
        : null
    } catch {
      return null
    }
  },
  set(val: CalendarDate | null) {
    projectState.value.endAt = val ? val.toString() : undefined
  }
})

const employeeStore = useEmployeeStore()
const { getEmployeeName } = employeeStore
const { employeeItems } = storeToRefs(employeeStore)

const isEditModalOpen = ref(false)
const isDeleteModalOpen = ref(false)
const isStartModalOpen = ref(false)
const isFinishModalOpen = ref(false)
const isResetModalOpen = ref(false)
const isCopyModalOpen = ref(false)
const selectedProject = ref<Project | null>(null)

const table = useTemplateRef("table")

const UDropdownMenu = resolveComponent("UDropdownMenu")
const UButton = resolveComponent("UButton")
const UBadge = resolveComponent("UBadge")
const UCheckbox = resolveComponent("UCheckbox")

// 刷新项目中的流程
function invalidProjectProcesses() {
  queryClient.invalidateQueries({
    queryKey: ["processes", selectedProject.value?.id]
  })
}

// 删除项目
async function handleDelete() {
  if (selectedProject.value) {
    try {
      await deleteProject(selectedProject.value.id)
      toast.add({
        title: `项目 ${selectedProject.value.name} 删除成功`,
        color: "success"
      })
    } catch (error) {
      toast.add({
        title: "删除失败",
        color: "error"
      })
      console.log("删除项目失败:", error)
    } finally {
      isDeleteModalOpen.value = false
    }
  }
}

// 开始项目
async function handleStart() {
  if (selectedProject.value) {
    try {
      await startProject(selectedProject.value.id)
      toast.add({
        title: `项目 ${selectedProject.value.name} 已开始`,
        color: "success"
      })
      invalidProjectProcesses()
    } catch (err) {
      toast.add({
        title: "开始失败",
        color: "error"
      })
      console.log("开始项目失败:", err)
    } finally {
      isStartModalOpen.value = false
    }
  }
}

// 完成项目
async function handleFinish() {
  if (selectedProject.value) {
    try {
      await finishProject(selectedProject.value.id)
      toast.add({
        title: `项目 ${selectedProject.value.name} 已完成`,
        color: "success"
      })
    } catch (err) {
      toast.add({
        title: "完成失败",
        color: "error"
      })
      console.log("完成项目失败:", err)
    } finally {
      isFinishModalOpen.value = false
    }
  }
}

// 重置项目
async function handleReset() {
  if (selectedProject.value) {
    try {
      await resetProject(selectedProject.value.id)
      toast.add({
        title: `项目 ${selectedProject.value.name} 已重置`,
        color: "success"
      })
      invalidProjectProcesses()
    } catch (err) {
      toast.add({
        title: "重置失败",
        color: "error"
      })
      console.log("重置项目失败:", err)
    } finally {
      isResetModalOpen.value = false
    }
  }
}

// query过滤
function handleSetFilter() {
  const cleanState = JSON.parse(JSON.stringify(projectState.value))
  setQuery({
    managerID: cleanState.managerID,
    status: cleanState.status,
    startAt: toRFC3339(cleanState.startAt),
    endAt: toRFC3339(cleanState.endAt)
  })
  setProjectState(cleanState)
  toast.add({
    title: "查询条件已更新",
    color: "success"
  })
}

// 重置query过滤条件
function handleResetFilter() {
  const defaultState = getDefaultProjectState()
  projectState.value = { ...defaultState }
  setQuery({
    managerID: defaultState.managerID,
    status: defaultState.status,
    startAt: toRFC3339(defaultState.startAt),
    endAt: toRFC3339(defaultState.endAt)
  })
  resetProjectState()
  toast.add({
    title: "查询条件已重置",
    color: "success"
  })
}

// 操作菜单定义
function getRowItems(row: Row<Project>) {
  return [
    {
      type: "label",
      label: "操作"
    },
    {
      type: "separator"
    },
    {
      label: "查看流程",
      icon: "i-lucide-eye",
      onSelect: () => {
        router.push(`/project/${row.original.id}`)
      }
    },
    {
      type: "separator"
    },
    {
      label: "编辑",
      icon: "i-lucide-pencil",
      disabled: !hasAtLeastRole("manager"),
      onSelect: () => {
        isEditModalOpen.value = true
        selectedProject.value = row.original
      }
    },
    {
      label: "开始",
      icon: "i-lucide-play",
      disabled: row.original.status !== "排单" || !hasAtLeastRole("manager"),
      onSelect: () => {
        isStartModalOpen.value = true
        selectedProject.value = row.original
      }
    },
    {
      label: "完成",
      icon: "i-lucide-check",
      disabled: row.original.status !== "进行" || !hasAtLeastRole("manager"),
      onSelect: () => {
        isFinishModalOpen.value = true
        selectedProject.value = row.original
      }
    },
    {
      label: "复制",
      icon: "i-lucide-copy",
      disabled: !hasAtLeastRole("manager"),
      onSelect: () => {
        isCopyModalOpen.value = true
        selectedProject.value = row.original
      }
    },
    {
      label: "导出",
      icon: "i-lucide-download",
      onSelect: () => {
        exportProject(row.original.id, row.original.name, toast)
      }
    },
    {
      type: "separator"
    },
    {
      label: "删除",
      icon: "i-lucide-trash",
      color: "error",
      disabled: !hasAtLeastRole("manager"),
      onSelect: () => {
        isDeleteModalOpen.value = true
        selectedProject.value = row.original
      }
    },
    {
      label: "重置",
      icon: "i-lucide-rotate-ccw",
      color: "error",
      disabled: !hasAtLeastRole("manager"),
      onSelect: () => {
        isResetModalOpen.value = true
        selectedProject.value = row.original
      }
    }
  ]
}

const rowSelection = ref<Record<number, boolean>>({})

// 表格列定义
const columns: TableColumn<Project>[] = [
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
    accessorKey: "id",
    header: "ID"
  },
  {
    accessorKey: "name",
    header: "名称",
    cell: ({ row }) => {
      return h(
        "div",
        {
          class: "cursor-pointer text-blue-600 hover:underline",
          onClick: () => router.push(`/project/${row.original.id}`)
        },
        row.original.name
      )
    }
  },
  {
    accessorKey: "description",
    header: "描述"
  },
  {
    header: "优先级",
    cell: ({ row }) => {
      const colorMap: Record<number, string> = {
        1: "success",
        2: "warning",
        3: "warning",
        4: "error",
        5: "error"
      }
      const color = colorMap[row.original.priority]
      return h(
        UBadge,
        {
          class: "capitalize",
          variant: "subtle",
          color
        },
        () => row.original.priority.toString()
      )
    }
  },
  {
    id: "manager",
    header: "负责人",
    accessorFn: (row) => row.managerID,
    filterFn: "equals",
    cell: ({ row }) => getEmployeeName(row.original.managerID)
  },
  {
    accessorKey: "startedAt",
    header: "开始时间",
    cell: ({ row }) => {
      const dateStr = row.original.startedAt
        ? row.original.startedAt.slice(0, 10)
        : ""

      const days = row.original.remainingDays
      const daysStr =
        days !== undefined
          ? days < 0
            ? ` (超时${-days}天)`
            : ` (剩余${days}天)`
          : ""

      return dateStr + daysStr
    }
  },
  {
    accessorKey: "status",
    header: "状态",
    filterFn: (row, columnID, filterValue: ProjectStatus[]) => {
      if (!filterValue.length) return false
      return filterValue.includes(row.getValue(columnID))
    },
    cell: ({ row }) => {
      const colorMap: Record<ProjectStatus, string> = {
        排单: "neutral",
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
    accessorKey: "remark",
    header: "备注"
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

// 项目名称过滤
const columnFilters = ref([
  {
    id: "name",
    value: ""
  }
])

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
