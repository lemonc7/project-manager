<template>
  <DropdownMenu>
    <DropdownMenuTrigger as-child>
      <Button variant="ghost" class="h-8 w-8 p-0">
        <MoreHorizontal />
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent>
      <DropdownMenuLabel>操作</DropdownMenuLabel>
      <DropdownMenuSeparator />
      <DropdownMenuItem @click="props.onEdit?.(props.process)">
        编辑
      </DropdownMenuItem>
      <DropdownMenuItem
        class="text-red-600 focus:text-red-600"
        @click="props.onDelete?.(props.process)"
      >
        删除
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>

<script setup lang="ts">
import type { Process } from "@/models/processes"
import { MoreHorizontal } from "lucide-vue-next"

const props = defineProps<{
  process: Process
  onEdit?: (process: Process) => void
  onDelete?: (process: Process) => void
}>()
</script>
