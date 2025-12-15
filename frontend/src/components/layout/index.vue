<template>
  <UApp
    :toaster="{
      duration: 2000,
      position: 'top-center'
    }"
  >
    <UDashboardGroup unit="rem" storage="local">
      <!-- 侧栏 -->
      <UDashboardSidebar
        id="default"
        v-model:open="open"
        collapsible
        resizable
        class="bg-elevated/25"
        :ui="{ footer: 'lg:border-t lg:border-default' }"
      >
        <template #header="{ collapsed }">
          <UButton
            variant="ghost"
            icon="material-symbols:folder-managed"
            block
            size="xl"
            color="neutral"
            :label="collapsed ? undefined : '项目管理平台'"
            :square="collapsed"
            @click="router.push('/')"
          />
        </template>

        <template #default="{ collapsed }">
          <UNavigationMenu
            :collapsed="collapsed"
            :items="links"
            orientation="vertical"
            tooltip
            popover
          />
        </template>
        <template #footer="{ collapsed }">
          <UserMenu :collapsed="collapsed" />
        </template>
      </UDashboardSidebar>
      <RouterView />
    </UDashboardGroup>
  </UApp>
</template>

<script setup lang="ts">
import type { NavigationMenuItem } from "@nuxt/ui"
import { ref, watch } from "vue"
import { useRouter } from "vue-router"
import { useEmployees } from "../../composables/useEmployees"
import { useTeams } from "../../composables/useTeams"
import { useEmployeeStore } from "../../stores/employee"
import { useTeamStore } from "../../stores/team"
import { useProcesses } from "../../composables/useProcesses"
import { useProcessStore } from "../../stores/process"

const open = ref(false)
const router = useRouter()

const links: NavigationMenuItem[] = [
  {
    label: "任务",
    to: "/",
    icon: "streamline-plump:task-list-edit",
    exact: true,
    onSelect: () => {
      open.value = false
    }
  },
  {
    label: "项目总览",
    to: "/project",
    icon: "i-lucide-sheet",
    exact: true,
    onSelect: () => {
      open.value = false
    }
  },
  {
    label: "团队",
    to: "/employee",
    defaultOpen: true,
    type: "trigger",
    icon: "fluent:organization-12-filled",
    children: [
      {
        label: "成员",
        to: "/employee",
        icon: "i-lucide-user",
        onSelect: () => {
          open.value = false
        }
      },
      {
        label: "小组",
        to: "/team",
        icon: "material-symbols:family-group",
        onSelect: () => {
          open.value = false
        }
      }
    ]
  }
]

const { employees } = useEmployees()
const { teams } = useTeams()
const { processes } = useProcesses()
const { setEmployees } = useEmployeeStore()
const { setTeams } = useTeamStore()
const { setMaxSequence } = useProcessStore()

// 同步employees
watch(
  () => employees.value,
  (list) => {
    if (list) {
      setEmployees(list)
    }
  },
  { immediate: true }
)

// 同步teams
watch(
  () => teams.value,
  (list) => {
    if (list) {
      setTeams(list)
    }
  },
  { immediate: true }
)

// 同步processes
watch(
  () => processes.value,
  (list) => {
    if (list) {
      setMaxSequence(list)
    }
  },
  { immediate: true }
)
</script>
