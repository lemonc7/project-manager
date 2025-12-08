<template>
  <UDashboardPanel id="task">
    <template #header>
      <UDashboardNavbar title="任务总览" :ui="{ right: 'gap-3' }">
        <template #leading>
          <UDashboardSidebarCollapse />
        </template>
        <template #right>
          <RefreshButton />
        </template>
      </UDashboardNavbar>
    </template>
    <template #body>
      <UPageGrid class="lg:grid-cols-3 gap-4 sm:gap-6 lg:gap-px">
        <UPageCard
          v-for="(stat, index) in baseStats"
          :key="index"
          :icon="stat.icon"
          :title="stat.title"
          :description="stat.description"
          variant="subtle"
          :ui="{
            container: 'gap-y-1.5',
            wrapper: 'items-start',
            leading:
              'p-2.5 rounded-full bg-primary/10 ring ring-inset ring-primary/25',
            title: 'font-normal text-muted text-xs uppercase'
          }"
          class="lg:rounded-none first:rounded-l-lg last:rounded-r-lg hover:z-1"
        >
          <div class="font-bold text-4xl">{{ dataLens[index] }}</div>
        </UPageCard>
      </UPageGrid>
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <TaskTable
          status="doing"
          :is-remaining-days-visible="true"
          @update-length="updateLen(0, $event)"
        />
        <TaskTable
          status="ready"
          class="hidden sm:block"
          :is-remaining-days-visible="false"
          @update-length="updateLen(1, $event)"
        />
        <TaskTable
          status="todo"
          class="hidden sm:block"
          :is-remaining-days-visible="false"
          @update-length="updateLen(2, $event)"
        />
      </div>
    </template>
  </UDashboardPanel>
</template>

<script setup lang="ts">
import { computed } from "vue"
import { useTaskStore } from "../composables/useTaskStore"

const baseStats = [
  {
    title: "进行",
    icon: "i-lucide-pause",
    description: "进行中的任务"
  },
  {
    title: "准备",
    icon: "i-lucide-play",
    description: "马上开始的任务"
  },
  {
    title: "完成",
    icon: "i-lucide-list-todo",
    description: "后续要做的任务"
  }
]

const taskStore = useTaskStore()

const dataLens = computed(() => taskStore.lengths)

function updateLen(index: number, length: number) {
  taskStore.setLength(index, length)
}
</script>
