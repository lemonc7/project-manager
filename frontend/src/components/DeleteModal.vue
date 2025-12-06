<template>
  <UModal
    v-model:open="props.open"
    @update:open="onOpenChange"
    :title="props.title"
    :description="props.description"
  >
    <template #body>
      <div class="flex justify-end gap-2">
        <UButton label="取消" variant="outline" @click="onOpenChange(false)" />
        <UButton
          label="删除"
          color="error"
          loading-auto
          @click="props.onSubmit"
        />
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
const props = defineProps<{
  open: boolean
  onSubmit: () => void
  title: string
  description: string
}>()

const emit = defineEmits<{
  (e: "update:open", value: boolean): void
}>()

function onOpenChange(open: boolean) {
  emit("update:open", open)
}
</script>
