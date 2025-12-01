<template>
  <span
    class="px-3 py-1 rounded-md text-white text-sm font-medium transition-all duration-200"
    :class="[
      statusClass,
      isClickable
        ? 'cursor-pointer hover:ring-2 hover:ring-blue-500 hover:ring-opacity-40 ring-offset-2 ring-offset-white'
        : ''
    ]"
    @click="handleClick"
  >
    {{ props.status }}
  </span>
</template>

<script setup lang="ts">
import type { ProcessStatus } from "@/models/processes"
import { computed } from "vue"

const props = defineProps<{
  status: ProcessStatus
}>()

const emit = defineEmits(["click"])

const statusColors: Record<string, string> = {
  等待: "bg-gray-500",
  准备: "bg-cyan-500",
  进行: "bg-blue-500",
  完成: "bg-green-500",
  暂停: "bg-yellow-500 text-black",
  终止: "bg-red-500"
}

const isClickable = computed(() => props.status === "进行")

const statusClass = computed(() => statusColors[props.status] || "bg-gray-300")

const handleClick = () => {
  if (isClickable.value) {
    emit("click")
  }
}
</script>
