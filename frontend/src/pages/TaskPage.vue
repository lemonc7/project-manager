<template>
  <UDashboardPanel id="task">
    <template #header>
      <UDashboardNavbar title="任务总览" :ui="{ right: 'gap-3' }">
        <template #leading>
          <UDashboardSidebarCollapse />
        </template>
        <template #right>
          <USelectMenu
            v-model="employeeID"
            :items="employeeItems"
            value-key="value"
            :search-input="{
              placeholder: '查看成员任务...'
            }"
            icon="i-lucide-user"
            :disabled="!hasAtLeastRole('manager')"
            class="w-32"
          />
          <RefreshButton />
        </template>
      </UDashboardNavbar>
    </template>
    <template #body>
      <div class="flex flex-col h-full gap-4">
        <div class="grid grid-cols-3">
          <UPageCard
            v-for="(stat, index) in baseStats"
            :key="index"
            :icon="stat.icon"
            :title="stat.title"
            :description="stat.description"
            variant="subtle"
            :ui="{
              leading:
                'p-2.5 rounded-full bg-info/10 ring ring-inset ring-info/25',
              title: 'font-normal text-muted text-xs'
            }"
            class="rounded-none first:rounded-l-lg last:rounded-r-lg transition-shadow duration-300 hover:shadow-lg"
          >
            <div class="font-bold text-4xl">{{ stat.length }}</div>
          </UPageCard>
        </div>
        <div
          class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 grow min-h-0"
        >
          <TaskTable
            :employee="employeeID"
            status="doing"
            :is-remaining-days-visible="true"
            @update-length="updateLen(0, $event)"
          />
          <TaskTable
            :employee="employeeID"
            status="ready"
            class="hidden md:flex"
            :is-remaining-days-visible="false"
            @update-length="updateLen(1, $event)"
          />
          <TaskTable
            :employee="employeeID"
            status="todo"
            class="hidden lg:flex"
            :is-remaining-days-visible="false"
            @update-length="updateLen(2, $event)"
          />
        </div>
      </div>
    </template>
  </UDashboardPanel>
</template>

<script setup lang="ts">
import { computed, ref } from "vue"
import { useTasks } from "../composables/useTasks"
import { getUserID, hasAtLeastRole } from "../utils/auth"
import { useEmployeeStore } from "../stores/employee"
import { storeToRefs } from "pinia"

const statTemplates = [
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

const { lengths, setLength } = useTasks()

const employeeID = ref(getUserID())
const employeeStore = useEmployeeStore()
const { employeeItems } = storeToRefs(employeeStore)

const baseStats = computed(() => {
  return statTemplates.map((t, i) => ({
    ...t,
    length: lengths.value[i] ?? 0
  }))
})

function updateLen(index: number, length: number) {
  setLength(index, length)
}
</script>
