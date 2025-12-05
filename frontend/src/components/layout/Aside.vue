<template>
  <UDashboardSidebar id="default" v-model:open="open" collapsible>
    
  </UDashboardSidebar>
</template>

<script setup lang="ts">
import type { NavigationMenuItem } from "@nuxt/ui"
import { ref } from "vue"

const open = ref(false)

const links: NavigationMenuItem[][] = [
  [
    {
      label: "任务",
      to: "/home",
      onSelect: () => {
        open.value = false
      }
    },
    {
      label: "项目总览",
      to: "/project",
      onSelect: () => {
        open.value = false
      }
    },
    {
      label: "团队",
      to: "/organization",
      defaultOpen: true,
      type: "trigger",
      children: [
        {
          label: "成员",
          to: "/organization/employee",
          onSelect: () => {
            open.value = false
          }
        },
        {
          label: "小组",
          to: "/organization/team",
          onSelect: () => {
            open.value = false
          }
        }
      ]
    }
  ]
]
</script>
