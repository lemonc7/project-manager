<template>
  <div
    v-if="isProjectsError || isEmployeesError"
    class="flex items-center justify-center h-full"
  >
    <span>加载失败</span>
  </div>
  <div
    v-else-if="isProjectsLoading || isEmployeesLoading"
    class="flex items-center justify-center h-full"
  >
    <span>项目列表加载中...</span>
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
</template>

<script setup lang="ts">
import EditableCell from "@/components/common/EditableCell.vue"
import useEmployees from "@/composables/use-employees"
import useProjects from "@/composables/use-projects"
import useUpdateProject from "@/composables/use-update-project"
import type {
  FilterStatus,
  ListProjectsRequest,
  Project
} from "@/models/projects"
import {
  FlexRender,
  getCoreRowModel,
  useVueTable,
  type ColumnDef
} from "@tanstack/vue-table"
import { h } from "vue"
import { toast } from "vue-sonner"
import ProjectActions from "./ProjectActions.vue"
import { useRouter } from "vue-router"
import ProjectStatus from "./ProjectStatus.vue"

const props = defineProps<{
  projectFilter: FilterStatus
  request: ListProjectsRequest
  disableSelect: boolean
}>()

const emit = defineEmits<{
  (e: "edit", project: Project): void
  (e: "delete", project: Project): void
}>()

const router = useRouter()

const {
  data: projects,
  isLoading: isProjectsLoading,
  isError: isProjectsError
} = useProjects({
  filter: props.projectFilter,
  req: props.request
})

const {
  data: employees,
  isLoading: isEmployeesLoading,
  isError: isEmployeesError
} = useEmployees()

// 更新项目
const updateProject = useUpdateProject({
  onError: () => {
    toast.error("更新失败")
  }
})

// 表格定义
const columns: ColumnDef<Project>[] = [
  {
    header: "项目名称",
    cell: ({ row }) => {
      const project = row.original
      return h(
        "div",
        {
          class: "cursor-pointer text-blue-600 hover:underline",
          onClick: () => router.push(`/project/${project.id}`)
        },
        project.name
      )
    }
  },
  {
    header: "细节描述",
    cell: ({ row }) => {
      return h(EditableCell, {
        value: row.original.description,
        onSave: (value: string) =>
          updateProject.mutate({
            id: row.original.id,
            req: { description: value }
          })
      })
    }
  },
  {
    header: "优先级",
    cell: ({ row }) => {
      const priority = row.original.priority
      if (priority <= 3) {
        return "低"
      }
      if (priority <= 6) {
        return "中"
      }
      return "高"
    }
  },
  {
    header: "负责人",
    cell: ({ row }) => {
      return h(EditableCell, {
        value: row.original.managerID,
        type: "select",
        options: employees.value?.map((e) => ({
          value: e.id,
          label: e.name
        })),
        onSave: (value: number) => {
          updateProject.mutate({
            id: row.original.id,
            req: { managerID: value }
          })
        },
        disabled: props.disableSelect
      })
    }
  },
  {
    header: "开始时间",
    cell: ({ row }) => {
      let value = ""
      if (row.original.startedAt) {
        value = new Date(row.original.startedAt).toISOString().split("T")[0]
      }
      if (row.original.remainingDays) {
        if (row.original.remainingDays < 0) {
          value += ` (超时${-row.original.remainingDays}天)`
        } else {
          value += ` (剩余${row.original.remainingDays}天)`
        }
      }
      return value
    }
  },
  {
    header: "状态",
    cell: ({ row }) => {
      return h(ProjectStatus, {
        status: row.original.status
      })
    }
  },
  {
    header: "备注",
    cell: ({ row }) => {
      return h(EditableCell, {
        value: row.original.remark,
        onSave: (value: string) =>
          updateProject.mutate({
            id: row.original.id,
            req: { remark: value }
          })
      })
    }
  },
  {
    id: "actions",
    size: 30,
    header: "操作",
    cell: ({ row }) => {
      const project = row.original
      return h(ProjectActions, {
        project,
        onEdit: () => emit("edit", project),
        onDelete: () => emit("delete", project)
      })
    }
  }
]

const table = useVueTable({
  get data() {
    return projects.value ?? []
  },
  columns,
  getCoreRowModel: getCoreRowModel()
})
</script>
