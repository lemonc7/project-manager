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
        <UInput
          v-model="name"
          class="max-w-sm"
          icon="i-lucide-search"
          placeholder="过滤项目名称..."
        />
        <div class="flex flex-wrap items-center gap-2">
          <UFieldGroup class="px-2">
            <USelect
              v-model="statusFilter"
              multiple
              placeholder="过滤状态..."
              icon="i-lucide-list-todo"
              :items="[...ProjectStatusOptions]"
              :ui="{
                trailingIcon:
                  'group-date-[state=open]:rotate-180 transition-transform duration-200'
              }"
              class="min-w-44"
            />
            <UTooltip :delay-duration="100" text="重置过滤条件">
              <UButton
                icon="i-lucide-list-restart"
                variant="outline"
                @click="statusFilter = ['排单', '进行', '暂停']"
              />
            </UTooltip>
          </UFieldGroup>

          <USelectMenu
            v-model="managerFilter"
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
      </div>
      <!-- table渲染 -->
      <UTable
        ref="table"
        v-model:column-filters="columnFilters"
        v-model:pagination="pagination"
        :pagination-options="{
          getPaginationRowModel: getPaginationRowModel()
        }"
        class="shrink-0"
        :data="projects ?? []"
        :columns="columns"
        :loading="isProjectLoading"
        :ui="{
          base: 'table-fixed border-separate border-spacing-0',
          thead: '[&>tr]:bg-elevated/50 [&>tr]:after:content-none',
          tbody: '[&>tr]:last:[&>td]:border-b-0',
          th: 'py-2 first:rounded-l-lg last:rounded-r-lg border-y border-default first:border-l last:border-r',
          td: 'py-2.5 border-b border-default',
          separator: 'h-0'
        }"
      />
      <div class="flex mt-auto justify-end">
        <UPagination
          :default-page="
            (table?.tableApi.getState().pagination.pageIndex || 0) + 1
          "
          :items-per-page="table?.tableApi.getState().pagination.pageSize"
          :total="totalRows"
          @update:page="(p) => table?.tableApi.setPageIndex(p - 1)"
        />
      </div>
    </template>
  </UDashboardPanel>
  <!-- 删除弹窗 -->
  <DeleteModal
    v-model:open="isDeleteModalOpen"
    :on-submit="handleDelete"
    title="删除项目"
    :description="`确认要删除[${selectedProject?.name}]项目么？该操作还会删除相关流程，无法撤销！`"
  />
</template>

<script setup lang="ts">
import { computed, h, ref, resolveComponent, useTemplateRef, watch } from "vue"
import useProjects from "../composables/useProjects"
import {
  getPaginationRowModel,
  type Row,
  type Table
} from "@tanstack/table-core"
import {
  ProjectStatusOptions,
  type Project,
  type ProjectStatus
} from "../models/projects"
import { useRouter } from "vue-router"
import type { TableColumn } from "@nuxt/ui"
import useEmployeeMap from "../composables/useEmployeeMap"
import useDeleteProject from "../composables/useDeleteProject"

const router = useRouter()
const toast = useToast()

const { data: projects, isLoading: isProjectLoading } = useProjects({
  filter: "filter",
  req: {}
})

const { employeeItems, getEmployeeName } = useEmployeeMap()

const isEditModalOpen = ref(false)
const isDeleteModalOpen = ref(false)
const selectedProject = ref<Project | null>(null)

const table = useTemplateRef<{ tableApi: Table<Project> }>("table")

const UDropdownMenu = resolveComponent("UDropdownMenu")
const UButton = resolveComponent("UButton")
const UBadge = resolveComponent("UBadge")

// 删除项目
const deleteProject = useDeleteProject({
  onSuccess: () => {
    toast.add({
      title: `项目 ${selectedProject.value?.name} 删除成功`,
      color: "success"
    })
    isDeleteModalOpen.value = false
  },
  onError: () => {
    toast.add({
      title: "删除失败",
      color: "error"
    })
  }
})

function handleDelete() {
  if (selectedProject.value) {
    deleteProject.mutate(selectedProject.value.id)
    isDeleteModalOpen.value = false
  }
}

// 分页定义
const pagination = ref({
  pageIndex: 0,
  pageSize: 10
})

const totalRows = computed(() => {
  return table.value?.tableApi.getFilteredRowModel().rows.length ?? 0
})

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
      label: "编辑",
      icon: "i-lucide-pencil",
      onSelect: () => {
        isEditModalOpen.value = true
        selectedProject.value = row.original
      }
    },
    {
      label: "删除",
      icon: "i-lucide-trash",
      color: "error",
      onSelect: () => {
        isDeleteModalOpen.value = true
        selectedProject.value = row.original
      }
    }
  ]
}

// 表格列定义
const columns: TableColumn<Project>[] = [
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
    header: "细节描述"
  },
  {
    header: "优先级",
    cell: ({ row }) => {
      const priority = row.original.priority
      const color =
        priority <= 3 ? "success" : priority <= 6 ? "warning" : "error"
      return h(
        UBadge,
        {
          class: "capitalize",
          variant: "subtle",
          color
        },
        () => priority.toString()
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
    header: "开始时间",
    cell: ({ row }) => {
      const dateStr = row.original.startedAt
        ? new Date(row.original.startedAt).toISOString().split("T")[0]
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

// status过滤
const statusFilter = ref<ProjectStatus[]>(["排单", "进行", "暂停"])
watch(statusFilter, (newVal) => {
  const statusColumn = table.value?.tableApi.getColumn("status")
  if (!statusColumn) return

  if (!newVal.length) {
    statusColumn.setFilterValue([])
  } else {
    statusColumn.setFilterValue(newVal)
  }
})

// 负责人过滤
const managerFilter = ref(0)
watch(managerFilter, (newVal) => {
  if (!table.value?.tableApi) return
  const managerColumn = table.value.tableApi.getColumn("manager")
  if (!managerColumn) return

  if (newVal === 0) {
    managerColumn.setFilterValue(undefined)
  } else {
    managerColumn.setFilterValue(newVal)
  }
})

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
