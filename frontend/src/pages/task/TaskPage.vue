<template>
  <div class="flex flex-col gap-4 p-4 h-full">
    <div class="w-full">
      <Card class="hover:shadow-lg transition-shadow duration-300">
        <CardHeader>
          <CardTitle>概览</CardTitle>
        </CardHeader>
        <CardContent>内容</CardContent>
      </Card>
    </div>
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <TaskTable status="doing" />
      <TaskTable status="ready" class="hidden sm:flex" />
      <TaskTable status="todo" class="hidden sm:flex" />
    </div>
  </div>
</template>

<script setup lang="ts">
import TaskTable from "./TaskTable.vue"
</script>
