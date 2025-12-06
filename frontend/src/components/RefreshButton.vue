<template>
  <UButton
    icon="i-lucide-refresh-cw"
    variant="outline"
    title="刷新"
    @click="refresh"
  />
</template>

<script setup lang="ts">
import { useQueryClient } from "@tanstack/vue-query"

const toast = useToast()
const queryClient = useQueryClient()
async function refresh() {
  await queryClient.refetchQueries({ type: "active" })
  toast.add({
    title: "刷新成功",
    description: "页面数据已更新",
    color: "success"
  })
}
</script>
