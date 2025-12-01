<template>
  <!-- 下拉选择框 -->
  <Select
    v-if="type === 'select'"
    :model-value="tempValue"
    :disabled="disabled"
    @update:model-value="(v) => props.onSave?.(v as any)"
  >
    <SelectTrigger>
      <SelectValue placeholder="请选择" />
    </SelectTrigger>
    <SelectContent>
      <SelectItem v-for="opt in options" :key="opt.value" :value="opt.value">
        {{ opt.label }}
      </SelectItem>
    </SelectContent>
  </Select>

  <!-- 非编辑状态 -->
  <div v-else-if="!editing" class="cursor-pointer" @click="editing = true">
    {{
      value !== null && value !== undefined && value !== ""
        ? String(value)
        : "-"
    }}
  </div>

  <!-- 数字输入框 -->
  <Input
    v-else-if="type === 'number'"
    type="number"
    v-model="tempValue"
    v-focus
    @blur="handleSave"
    @keydown.enter="handleSave"
    @keydown.esc="handleCancel"
  />

  <!-- 文本输入框 -->
  <Input
    v-else
    v-model="tempValue"
    v-focus
    @blur="handleSave"
    @keydown.enter="handleSave"
    @keydown.esc="handleCancel"
  />
</template>

<script setup lang="ts">
import { ref, watch } from "vue"

type EditableType = "text" | "number" | "select"

interface Option<T> {
  label: string
  value: T
}

interface Props<T> {
  value: T
  type?: EditableType
  options?: Option<T>[]
  onSave?: (value: T) => void
  disabled?: boolean
}

const props = defineProps<Props<any>>()

const editing = ref(false)
const tempValue = ref(props.value)

// 当外部value改变时，同步tempValue
watch(
  () => props.value,
  (val) => {
    tempValue.value = val
  }
)

// 自动聚焦到input，兼容普通input和框架中的input
const vFocus = {
  mounted: (el: HTMLElement) => {
    const input = el.tagName === "INPUT" ? el : el.querySelector("input")
    input?.focus()
  }
}

function handleSave() {
  editing.value = false
  if (tempValue.value !== props.value) {
    props.onSave?.(tempValue.value)
  }
}

function handleCancel() {
  editing.value = false
  tempValue.value = props.value
}
</script>
