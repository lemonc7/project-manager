<template>
  <DropdownMenu>
    <DropdownMenuTrigger as-child>
      <Button variant="ghost" class="h-8 w-8 p-0">
        <MoreHorizontal />
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent>
      <DropdownMenuItem @click="props.onEdit?.(props.employee)">
        编辑
      </DropdownMenuItem>
      <DropdownMenuItem
        class="text-red-600 focus:text-red-600"
        @click="props.onDelete?.(props.employee)"
      >
        删除
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>

<script setup lang="ts">
import type { Employee } from "@/model/employees"
import { MoreHorizontal } from "lucide-vue-next"

const props = defineProps<{
  employee: Employee
  onEdit?: (employee: Employee) => void
  onDelete?: (employee: Employee) => void
}>()
</script>
