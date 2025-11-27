<template>
  <Card
    :class="['hover:shadow-lg transition-shadow duration-300', props.class]"
  >
    <CardHeader>
      <CardTitle>{{ titleMap[props.status] }}</CardTitle>
    </CardHeader>
    <CardContent>
      <div v-if="isLoading" class="flex items-center justify-center h-full">
        任务加载中...
      </div>
      <div v-else-if="isError" class="flex items-center justify-center h-full">
        加载失败
      </div>
      <Table v-else>
        <TableHeader>
          <TableRow>
            <TableHead>任务</TableHead>
            <TableHead>详情</TableHead>
            <TableHead>项目</TableHead>
            <TableHead v-if="props.status === 'doing'">剩余天数</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow
            v-for="t in data"
            :key="t.id"
            class="cursor-pointer"
            @click="goProject(t.projectID)"
          >
            <TableCell>{{ t.name }}</TableCell>
            <TableCell>{{ t.description }}</TableCell>
            <TableCell>{{ t.projectName }}</TableCell>
            <TableCell v-if="props.status === 'doing'">{{
              t.remainingDays
            }}</TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </CardContent>
  </Card>
</template>

<script setup lang="ts">
import { listTasks } from "@/api/processes"
import Card from "@/components/ui/card/Card.vue"
import type { TaskStatus } from "@/model/projects"
import { getUserID } from "@/util/auth"
import { useQuery } from "@tanstack/vue-query"
import { useRouter } from "vue-router"

const props = defineProps<{ status: TaskStatus; class?: string }>()

const router = useRouter()
const userID = getUserID()

const { data, isLoading, isError } = useQuery({
  queryKey: ["tasks", props.status, { userID }],
  queryFn: () => {
    if (!userID) {
      throw new Error("缺少用户ID")
    }
    return listTasks(+userID, props.status)
  }
})

const titleMap: Record<TaskStatus, string> = {
  todo: "待办",
  ready: "准备",
  doing: "进行"
}

function goProject(id: number) {
  router.push(`/project/${id}`)
}
</script>
