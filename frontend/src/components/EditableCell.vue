<template>
  <!-- 非编辑状态 -->
  <span v-show="!editing" class="cursor-pointer" @click="startEdit">
    {{
      value !== null && value !== undefined && value !== ""
        ? String(value)
        : "-"
    }}
  </span>

  <!-- 编辑状态 -->
  <input
    v-show="editing"
    ref="inputRef"
    :type="type"
    v-model="tempValue"
    class="border-2 border-black focus:outline-none rounded px-2 py-1 box-border"
    @blur="handleSave"
    @keydown.enter="handleSave"
    @keydown.esc="handleCancel"
  />
</template>

<script setup lang="ts">
import { nextTick, ref, watch } from "vue"

type EditableType = "text" | "number"

interface Props<T> {
  value: T
  type?: EditableType
  onSave?: (value: T) => void
}

const props = defineProps<Props<any>>()

const editing = ref(false)
const tempValue = ref(props.value)

// 当外部value改变时，同步tempValue
watch(
  () => props.value,
  (val) => {
    if (!editing.value) {
      tempValue.value = val
    }
  }
)

const inputRef = ref<HTMLInputElement | null>(null)

function startEdit() {
  editing.value = true
  nextTick(() => {
    // 手动聚焦并选中
    inputRef.value?.focus()
    inputRef.value?.select?.()
  })
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
