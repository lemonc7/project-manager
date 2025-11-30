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
      <DropdownMenuItem @click="props.onEdit?.(props.project)">
        编辑
      </DropdownMenuItem>
      <DropdownMenuItem
        class="text-red-600 focus:text-red-600"
        @click="props.onDelete?.(props.project)"
      >
        删除
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>

<script setup lang="ts">
import type { Project } from "@/models/projects"
import { MoreHorizontal } from "lucide-vue-next";

const props = defineProps<{
  project: Project
  onEdit?: (project: Project) => void
  onDelete?: (project: Project) => void
}>()
</script>
