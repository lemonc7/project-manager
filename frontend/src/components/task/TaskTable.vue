<template>
  <div
    :class="[
      'flex flex-col h-full min-h-0 bg-white rounded-xl',
      'border border-default transition-shadow duration-300 hover:shadow-xl',
      props.class
    ]"
  >
    <div class="px-4 py-4 shrink-0">
      <UBadge
        size="xl"
        :color="titleMap[props.status].color"
        :icon="titleMap[props.status].icon"
        variant="subtle"
        >{{ titleMap[props.status].label }}</UBadge
      >
    </div>
    <UTable
      sticky="header"
      :data="data"
      :columns="columns"
      class="px-4"
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
</template>

<script setup lang="ts">
import { useRouter } from "vue-router"
import type { TaskStatus } from "../../models/projects"
import { useQuery } from "@tanstack/vue-query"
import { listTasks } from "../../api/processes"
import type { TableColumn } from "@nuxt/ui"
import type { Task } from "../../models/processes"
import { computed, h, watch } from "vue"
import { useProcessStore } from "../../stores/process"

const props = defineProps<{
  employee?: number
  status: TaskStatus
  class?: string
  isRemainingDaysVisible: boolean
}>()

const emit = defineEmits<{
  (e: "updateLength", length: number): void
}>()

const router = useRouter()
const { data, isLoading } = useQuery({
  queryKey: computed(() => [
    "tasks",
    props.status,
    { employeeID: props.employee }
  ]),
  queryFn: () => {
    if (!props.employee) {
      console.error("缺少用户ID")
      return
    }
    return listTasks(props.employee, props.status)
  },
  staleTime: Infinity,
  refetchInterval: 1000 * 60
})

const titleMap: Record<
  TaskStatus,
  { label: string; color: string; icon: string }
> = {
  todo: { label: "待办", color: "success", icon: "i-lucide-list-todo" },
  ready: { label: "准备", color: "warning", icon: "i-lucide-play" },
  doing: { label: "进行", color: "error", icon: "i-lucide-pause" }
}

const { setSequenceToScroll } = useProcessStore()

const columns: TableColumn<Task>[] = [
  {
    accessorKey: "name",
    header: "任务",
    cell: ({ row }) => {
      return h(
        "span",
        {
          style: { color: "blue", cursor: "pointer" },
          onClick: () =>
            goToProcess(row.original.sequence, row.original.projectID)
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

// 点击跳转
function goToProcess(sequence: number, projectID: number) {
  setSequenceToScroll(sequence)
  router.push(`/project/${projectID}`)
}

// 监听data变化，更新长度
watch(
  () => data.value?.length,
  (len) => {
    if (len !== undefined) {
      emit("updateLength", len)
    }
  }
)
</script>
