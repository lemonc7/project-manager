<template>
  <UDrawer
    :open="props.open"
    @update:open="onOpenChange"
    :title="`${props.type === 'rework' ? '返工' : '完成'}流程`"
    description="请输入机时和备注"
    :ui="{ container: 'max-w-xl mx-auto' }"
  >
    <template #body>
      <UForm
        ref="formRef"
        :schema="schema"
        :state="state"
        class="space-y-6"
        @submit="handleSubmit"
      >
        <!-- 机时 -->
        <UFormField label="机时" name="spendTime" required>
          <UInputNumber
            v-model="state.spendTime"
            class="w-full"
            placeholder="请输入机时"
            :min="0"
            :step="0.5"
          />
        </UFormField>
        <!-- 备注 -->
        <UFormField label="备注" name="remark">
          <UTextarea
            v-model="state.remark"
            class="w-full"
            placeholder="请输入备注"
            autoresize
          />
        </UFormField>
      </UForm>
    </template>

    <template #footer>
      <UButton
        label="确认"
        color="neutral"
        class="justify-center"
        :loading="props.type === 'rework' ? isReworking : isFinishing"
        @click="formRef.submit()"
      />
      <UButton
        label="取消"
        variant="outline"
        class="justify-center"
        @click="onOpenChange(false)"
      />
    </template>
  </UDrawer>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from "vue"
import z from "zod"
import { useProcessStore } from "../../composables/useProcessStore"
import { storeToRefs } from "pinia"
import type { FormSubmitEvent } from "@nuxt/ui"
import type { Process } from "../../models/processes"

const props = defineProps<{
  open: boolean
  type: "rework" | "finish"
  process: Process | null
}>()

const emit = defineEmits<{
  (e: "update:open", open: boolean): void
}>()

const toast = useToast()

const processStore = useProcessStore()
const { finishProcess, reworkProcess } = processStore
const { isFinishing, isReworking } = storeToRefs(processStore)
const schema = z.object({
  spendTime: z.number().min(0, "时间不能小于0"),
  remark: z.string().default("")
})
const formRef = ref()
type Schema = z.output<typeof schema>

const initialState: Schema = {
  spendTime: 0,
  remark: ""
}

const state = reactive<Schema>({ ...initialState })

function onOpenChange(open: boolean) {
  emit("update:open", open)
}

// 表单提交
async function handleSubmit(event: FormSubmitEvent<Schema>) {
  if (!props.process) return
  if (props.type === "rework") {
    try {
      await reworkProcess({
        id: props.process.id,
        req: {
          spendTime: event.data.spendTime,
          remark: event.data.remark
        }
      })
      toast.add({
        title: `流程 ${props.process.name} 已添加返工状态`,
        color: "success"
      })
      onOpenChange(false)
    } catch (err) {
      toast.add({
        title: "添加返工状态失败",
        color: "error"
      })
      console.log("添加返工状态失败:", err)
    }
  } else {
    try {
      await finishProcess({
        id: props.process.id,
        req: {
          spendTime: event.data.spendTime,
          remark: event.data.remark
        }
      })
      toast.add({
        title: `流程 ${props.process.name} 已完成`,
        color: "success"
      })
      onOpenChange(false)
    } catch (err) {
      toast.add({
        title: "流程完成设置失败",
        color: "error"
      })
      console.log("流程完成设置失败:", err)
    }
  }
}

// 监视抽屉框，打开时重置表单
watch(
  () => props.open,
  (isOpen) => {
    if (isOpen) {
      Object.assign(state, initialState)
    }
  }
)
</script>
