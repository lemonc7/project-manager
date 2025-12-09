<template>
  <UDashboardPanel id="team">
    <template #header>
      <UDashboardNavbar title="小组列表">
        <template #leading>
          <UDashboardSidebarCollapse title="折叠侧栏" />
        </template>
        <template #right>
          <div class="flex flex-wrap items-center gap-2">
            <CreateTeamModal />
            <RefreshButton />
          </div>
        </template>
      </UDashboardNavbar>
    </template>

    <template #body>
      <!-- table渲染 -->
      <UTable
        sticky
        ref="table"
        class="shrink-0"
        :loading="isLoading"
        :data="teams"
        :columns="columns"
        :ui="{
          base: 'table-fixed border-separate border-spacing-0',
          thead: '[&>tr]:bg-elevated/50 [&>tr]:after:content-none',
          tbody: '[&>tr]:last:[&>td]:border-b-0',
          th: 'py-2 first:rounded-l-lg last:rounded-r-lg border-y border-default first:border-l last:border-r',
          td: 'py-2 border-b border-default',
          separator: 'h-0'
        }"
      />
    </template>
  </UDashboardPanel>
  <!-- 删除弹窗 -->
  <AlertModal
    v-model:open="isDeleteModalOpen"
    :on-submit="handleDelete"
    title="删除分组"
    label="删除"
    color="error"
    :description="`确认要删除[${selectedTeam?.name}]么？该操作无法撤销!`"
  />
  <!-- 编辑弹窗 -->
  <EditTeamModal v-model:open="isEditModalOpen" :team="selectedTeam" />
</template>

<script setup lang="ts">
import type { TableColumn } from "@nuxt/ui"
import { h, ref, resolveComponent } from "vue"
import { type Row } from "@tanstack/table-core"
import { useTeamStore } from "../composables/useTeamStore"
import { storeToRefs } from "pinia"
import type { Team } from "../models/teams"

// 数据获取
const teamStore = useTeamStore()
const { teams, isLoading } = storeToRefs(teamStore)
const { deleteTeam } = teamStore

const toast = useToast()

// 模态框
const isEditModalOpen = ref(false)
const isDeleteModalOpen = ref(false)
const selectedTeam = ref<Team | null>(null)

async function handleDelete() {
  if (selectedTeam.value) {
    try {
      await deleteTeam(selectedTeam.value.id)
      toast.add({
        title: `分组 ${selectedTeam.value.name} 删除成功`,
        color: "success"
      })
    } catch (error) {
      toast.add({
        title: "删除失败",
        color: "error"
      })
      console.log("删除分组失败:", error)
    } finally {
      isDeleteModalOpen.value = false
    }
  }
}

// 组件定义
const UDropdownMenu = resolveComponent("UDropdownMenu")
const UButton = resolveComponent("UButton")

// 操作菜单定义
function getRowItems(row: Row<Team>) {
  return [
    {
      type: "label",
      label: "操作"
    },
    {
      type: "separator"
    },
    {
      label: "编辑",
      icon: "i-lucide-pencil",
      onSelect: () => {
        isEditModalOpen.value = true
        selectedTeam.value = row.original
      }
    },
    {
      label: "删除",
      icon: "i-lucide-trash",
      color: "error",
      onSelect: () => {
        isDeleteModalOpen.value = true
        selectedTeam.value = row.original
      }
    }
  ]
}

// table列定义
const columns: TableColumn<Team>[] = [
  {
    accessorKey: "id",
    header: "ID"
  },
  {
    accessorKey: "name",
    header: "小组"
  },

  {
    accessorKey: "description",
    header: "描述"
  },
  {
    id: "actions",
    cell: ({ row }) => {
      return h(
        "div",
        { class: "text-right" },
        h(
          UDropdownMenu,
          {
            content: {
              align: "end"
            },
            items: getRowItems(row)
          },
          () =>
            h(UButton, {
              icon: "i-lucide-ellipsis-vertical",
              variant: "ghost",
              class: "ml-auto"
            })
        )
      )
    }
  }
]
</script>
