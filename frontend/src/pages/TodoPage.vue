<template>
  <UDashboardPanel id="todo">
    <template #header>
      <UDashboardNavbar title="todo列表">
        <template #leading>
          <UDashboardSidebarCollapse title="折叠侧栏" />
        </template>
        <template #right>
          <div class="flex flex-wrap items-center gap-2">
            <RefreshButton />
          </div>
        </template>
      </UDashboardNavbar>
    </template>

    <template #body> </template>
  </UDashboardPanel>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { useTodos } from "../composables/useTodos"
import { getTodoState } from "../utils/todoState"
import { useEmployeeStore } from "../stores/employee"
import { storeToRefs } from "pinia"
import type { Todo } from "../models/todos"

defineOptions({
  name: "todo"
})

const toast = useToast()
const { todosReporter, todosAssignee, setQuery } = useTodos()

const todoState = ref(getTodoState())
const employeeStore = useEmployeeStore()
const { employeeItems } = storeToRefs(employeeStore)

const isEditModalOpen = ref(false)
const isDeleteModalOpen = ref(false)
const selectedTodo = ref<Todo | null>(null)

const tabItems = [
  {
    label: "我分配的",
    icon: "streamline-ultimate:task-finger-show-bold",
    slot: "reporter"
  },
  {
    label: "我负责的",
    icon: "i-lucide-list-todo",
    slot: "assignee"
  }
]
</script>
