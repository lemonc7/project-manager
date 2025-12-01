<template>
  <span
    class="px-3 py-1 rounded-md text-white text-sm font-medium"
    :class="statusClass"
    @click="handleClick"
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
const emit = defineEmits(["click"])

const statusColors: Record<string, string> = {
  排单: "bg-gray-500",
  进行: "bg-blue-500",
  完成: "bg-green-500",
  暂停: "bg-yellow-500 text-black",
  终止: "bg-red-500"
}

const isStartClickable = computed(() => props.status === "排单")
const isFinishClickable = computed(() => props.status === "进行")

const statusClass = computed(() => {
  const baseStatusClass = statusColors[props.status] || "bg-gray-300"

  if (isStartClickable.value || isFinishClickable.value) {
    let clickableClasses =
      "cursor-pointer hover:ring-2 hover:ring-opacity-40 ring-offset-2 ring-offset-white"

    if (isStartClickable.value) {
      clickableClasses += " hover:ring-gray-500"
    } else if (isFinishClickable.value) {
      clickableClasses += " hover:ring-blue-500"
    }

    return `${baseStatusClass} ${clickableClasses}`
  }

  return baseStatusClass
})

const handleClick = () => {
  if (isStartClickable.value) {
    emit("click", "start")
  } else if (isFinishClickable.value) {
    emit("click", "finish")
  }
}
</script>
