<template>
  <div
    :class="[
      'bg-white/50 dark:bg-gray-900/50 backdrop-blur border border-default rounded-xl shadow-md transition-shadow duration-300',
      'hover:shadow-lg',
      props.class
    ]"
  >
    <div class="px-4 py-3">
      <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
        {{ titleMap[props.status] }}
      </h2>
    </div>

    <div class="px-4 pt-2 pb-4">
      <UTable
        sticky
        :data="data"
        :columns="columns"
        class="shrink-0"
        :loading="isLoading"
        :column-visibility="{
          remainingDays: props.isRemainingDaysVisible
        }"
        :ui="{
          base: 'table-fixed border-separate border-spacing-0',
          thead: '[&>tr]:bg-elevated/50 [&>tr]:after:content-none',
          tbody: '[&>tr]:last:[&>td]:border-b-0',
          th: 'py-2 first:rounded-l-lg last:rounded-r-lg border-y border-default first:border-l last:border-r',
          td: 'py-2.5 border-b border-default',
          separator: 'h-0'
        }"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router"
import type { TaskStatus } from "../../models/projects"
import { getUserID } from "../../utils/auth"
import { useQuery } from "@tanstack/vue-query"
import { listTasks } from "../../api/processes"
import type { TableColumn } from "@nuxt/ui"
import type { Task } from "../../models/processes"
import { h, watch } from "vue"

const props = defineProps<{
  status: TaskStatus
  class?: string
  isRemainingDaysVisible: boolean
}>()

const emit = defineEmits<{
  (e: "updateLength", length: number): void
}>()

const router = useRouter()
const userID = getUserID()
const { data, isLoading } = useQuery({
  queryKey: ["tasks", props.status, { userID }],
  queryFn: () => {
    if (!userID) {
      throw new Error("缺少用户ID")
    }
    return listTasks(userID, props.status)
  }
})
const titleMap: Record<TaskStatus, string> = {
  todo: "待办",
  ready: "准备",
  doing: "进行"
}

const columns: TableColumn<Task>[] = [
  {
    accessorKey: "name",
    header: "任务",
    cell: ({ row }) => {
      return h(
        "span",
        {
          style: { color: "blue", cursor: "pointer" },
          onClick: () => router.push(`/project/${row.original.projectID}`)
        },
        row.original.name
      )
    }
  },
  {
    accessorKey: "description",
    header: "详情"
  },
  {
    accessorKey: "projectName",
    header: "项目"
  },
  {
    accessorKey: "remainingDays",
    header: "剩余时间",
    cell: ({ row }) => {
      if (row.original.remainingDays === undefined) {
        return "-"
      } else if (row.original.remainingDays >= 0) {
        return `剩余${row.original.remainingDays}天`
      } else {
        return `超时${-row.original.remainingDays}天`
      }
    }
  }
]

watch(
  () => data.value?.length,
  (len) => {
    if (len !== undefined) {
      emit("updateLength", len)
    }
  }
)
</script>
