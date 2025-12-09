<template>
  <UModal
    :open="props.open"
    @update:open="onOpenChange"
    :title="props.title"
    :description="props.description"
  >
    <template #body>
      <div class="flex justify-end gap-2">
        <UButton label="取消" variant="outline" @click="onOpenChange(false)" />
        <UButton
          :label="label"
          :color="props.color"
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
  color: string
  label: string
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
