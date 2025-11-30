<template>
  <span
    class="px-3 py-1 rounded-md text-white text-sm font-medium"
    :class="statusClass"
  >
    {{ props.status }}
  </span>
</template>

<script setup lang="ts">
import type { ProcessStatus } from "@/models/processes";
import { computed } from "vue"

const props = defineProps<{
  status: ProcessStatus
}>()

const statusColors: Record<string, string> = {
  等待: "bg-gray-500",
  准备: "bg-cyan-500",
  进行: "bg-blue-500",
  完成: "bg-green-500",
  暂停: "bg-yellow-500 text-black",
  终止: "bg-red-500"
}

const statusClass = computed(() => statusColors[props.status] || "bg-gray-300")
</script>
