<template>
  <UModal
    :open="props.open"
    @update:open="onOpenChange"
    title="编辑"
    :description="detail"
  >
    <template #body>
      <UForm
        :schema="schema"
        :state="state"
        class="space-y-4"
        @submit="handleSubmit"
      >
        <!-- 流程名称 -->
        <UFormField label="名称" name="name" required>
          <UInput
            v-model="state.name"
            class="w-full"
            placeholder="请输入流程名称"
          />
        </UFormField>
        <!-- 细节描述 -->
        <UFormField label="细节描述" name="description">
          <UInput
            v-model="state.description"
            class="w-full"
            placeholder="请输入细节描述"
          />
        </UFormField>
        <!-- 负责人+状态+机时 -->
        <div class="grid grid-cols-3 gap-4">
          <UFormField label="负责人" name="manager">
            <USelectMenu
              v-model="state.employee"
              :items="employeeItems"
              value-key="value"
              placeholder="请选择负责人"
              class="w-full"
            />
          </UFormField>
          <UFormField label="状态" name="status" required>
            <USelect
              v-model="state.status"
              :items="[...ProcessStatusOptions]"
              placeholder="请选择状态"
              class="w-full"
            />
          </UFormField>
          <UFormField label="机时" name="spendTime">
            <UInput
              type="number"
              v-model="state.spendTime"
              class="w-full"
              placeholder="请输入机时"
            />
          </UFormField>
        </div>
        <!-- 备注 -->
        <UFormField>
          <UTextarea
            v-model="state.remark"
            class="w-full"
            placeholder="请输入备注"
            autoresize
          />
        </UFormField>
        <!-- 按钮 -->
        <div class="flex justify-end gap-2">
          <UButton
            label="取消"
            variant="outline"
            @click="onOpenChange(false)"
          />
          <UButton label="更新" type="submit" :loading="isUpdating" />
        </div>
      </UForm>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import z from "zod"
import { computed, reactive, ref, watch } from "vue"
import { storeToRefs } from "pinia"
import { useEmployeeStore } from "../../composables/useEmployeeStore"
import {
  ProcessStatusOptions,
  type Process,
  type UpdateProcessRequest
} from "../../models/processes"
import { useProcessStore } from "../../composables/useProcessStore"

const props = defineProps<{
  open: boolean
  process: Process | null
}>()

const emit = defineEmits<{
  (e: "update:open", open: boolean): void
}>()

const toast = useToast()

// 定义表单数据
const schema = z.object({
  name: z.string().nonempty("请输入姓名"),
  description: z.string().optional(),
  remark: z.string().optional(),
  status: z.enum(["等待", "准备", "进行", "完成", "暂停", "终止"]),
  employee: z.number().optional(),
  spendTime: z.number().min(0, "花费时间不能小于0").optional()
})

type Schema = z.output<typeof schema>

const state = reactive<Schema>({
  name: props.process?.name ?? "",
  description: props.process?.description,
  remark: props.process?.remark,
  status: props.process?.status ?? "等待",
  employee: props.process?.employeeID,
  spendTime: props.process?.spendTime
})

function onOpenChange(open: boolean) {
  emit("update:open", open)
}

const processStore = useProcessStore()
const employeeStore = useEmployeeStore()
const { updateProcess } = processStore
const { isUpdating } = storeToRefs(processStore)
const { employeeItems } = storeToRefs(employeeStore)

async function handleSubmit() {
  if (!props.process?.id) return
  const req: UpdateProcessRequest = {
    name: state.name,
    description: state.description,
    remark: state.remark,
    status: state.status,
    employeeID: state.employee,
    spendTime: state.spendTime
  }
  try {
    await updateProcess({ id: props.process.id, req })
    toast.add({
      title: `流程 ${state.name} 更新成功`,
      color: "success"
    })
    onOpenChange(false)
  } catch (err) {
    toast.add({
      title: "更新失败",
      color: "error"
    })
    console.log("更新流程失败:", err)
  }
}

const detail = computed(() => {
  if (!props.process?.startedAt) return "流程未开始"
  const text = props.process.startedAt.slice(0, 10)
  if (props.process.spendTime !== undefined) {
    return `流程开始时间:${text}(机时${props.process.spendTime}h)`
  }
  return `流程开始时间:${text}`
})

// 监听open变化，打开时同步表单数据
watch(
  () => props.open,
  (isOpen) => {
    if (isOpen && props.process) {
      Object.assign(state, {
        name: props.process.name,
        description: props.process.description,
        remark: props.process.remark,
        status: props.process.status,
        employee: props.process.employeeID,
        spendTime: props.process.spendTime
      })
    }
  }
)
</script>
