<template>
  <span
    class="px-3 py-1 rounded-md text-white text-sm font-medium"
    :class="statusClass"
  >
    {{ props.status }}
  </span>
</template>

<script setup lang="ts">
import type { ProjectStatus } from "@/models/projects"
import { computed } from "vue"

const props = defineProps<{
  status: ProjectStatus
}>()

const statusColors: Record<string, string> = {
  排单: "bg-gray-500",
  进行: "bg-blue-500",
  完成: "bg-green-500",
  暂停: "bg-yellow-500 text-black",
  终止: "bg-red-500"
}

const statusClass = computed(() => statusColors[props.status] || "bg-gray-300")
</script>
