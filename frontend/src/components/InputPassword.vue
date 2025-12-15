<template>
  <UInput
    :type="show ? 'text' : 'password'"
    :value="props.modelValue"
    class="w-full"
    :placeholder="props.placeholder"
    :ui="{ trailing: 'pe-1' }"
    @input="emit('update:modelValue', $event.target.value)"
  >
    <template #trailing>
      <UButton
        variant="link"
        size="sm"
        :icon="show ? 'i-lucide-eye-off' : 'i-lucide-eye'"
        :aria-label="show ? '隐藏密码' : '显示密码'"
        :aria-pressed="show"
        aria-controls="password"
        @click="show = !show"
      />
    </template>
  </UInput>
</template>

<script setup lang="ts">
import { ref } from "vue"

const props = defineProps<{
  modelValue: string
  placeholder?: string
}>()

const emit = defineEmits<{
  (e: "update:modelValue", value: string): void
}>()

const show = ref(false)
</script>
