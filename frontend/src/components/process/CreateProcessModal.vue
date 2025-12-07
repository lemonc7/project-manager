<template>
  <UModal
    v-model:open="open"
    title="创建"
    description="为当前项目添加一个新流程"
  >
    <UButton label="创建流程" icon="i-lucide-git-pull-request-create-arrow" />
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
            placeholder="请输入项目名称"
            autofocus
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
        <!-- 负责人+流程序号 -->
        <div class="grid grid-cols-2 gap-4">
          <UFormField label="负责人" name="employee">
            <USelectMenu
              v-model="state.employee"
              :items="employeeItems"
              value-key="value"
              placeholder="请选择负责人"
              class="w-full"
            />
          </UFormField>
          <UFormField label="流程序号" name="sequence">
            <UInput
              type="number"
              v-model="state.sequence"
              class="w-full"
              placeholder="请输入流程序号"
            />
          </UFormField>
        </div>
        <!-- 备注 -->
        <UFormField label="备注" name="remark">
          <UTextarea v-model="state.remark" autoresize class="w-full" />
        </UFormField>
        <!-- 按钮 -->
        <div class="flex justify-end gap-2">
          <UButton label="取消" variant="outline" @click="open = false" />
          <UButton
            label="创建"
            type="submit"
            :loading="createProcess.isPending.value"
          />
        </div>
      </UForm>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import type { FormSubmitEvent } from "@nuxt/ui"
import { reactive, ref, watch } from "vue"
import z from "zod"
import { useProcesses } from "../../composables/useProcesses"
import useCreateProcess from "../../composables/useCreateProcess"
import type { InsertProcessRequest } from "../../models/processes"
import useEmployeeMap from "../../composables/useEmployeeMap"

const open = ref(false)
const { maxSequence } = useProcesses()
const { employeeItems } = useEmployeeMap()

// 定义表单数据
const schema = z.object({
  name: z.string().nonempty("请输入流程名称"),
  sequence: z
    .number()
    .min(1, "序号不能小于1")
    .max(maxSequence.value, `序号不能大于${maxSequence.value}`)
    .optional(),
  employee: z.number().optional(),
  description: z.string().optional(),
  remark: z.string().optional()
})

type Schema = z.output<typeof schema>

// 表单默认数据状态
const initialState: Schema = {
  name: "",
  sequence: undefined,
  employee: undefined,
  description: undefined,
  remark: undefined
}
const state = reactive<Schema>({ ...initialState })

const toast = useToast()
// 创建流程
const createProcess = useCreateProcess({
  onSuccess: () => {
    toast.add({
      title: `流程 ${state.name} 创建成功`,
      color: "success"
    })
    open.value = false
  },
  onError: () => {
    toast.add({
      title: "创建失败",
      color: "error"
    })
  }
})

// 表单提交
function handleSubmit(event: FormSubmitEvent<Schema>) {
  const req: Omit<InsertProcessRequest, "projectID"> = {
    name: event.data.name,
    employeeID: event.data.employee,
    description: event.data.description,
    remark: event.data.remark
  }

  createProcess.mutate({ sequence: event.data.sequence, req })
}

// 监视弹窗，关闭时重置表单
watch(open, (newVal) => {
  if (!newVal) {
    Object.assign(state, initialState)
  }
})
</script>
