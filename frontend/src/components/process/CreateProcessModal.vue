<template>
  <UModal
    v-model:open="open"
    title="创建"
    description="为当前项目添加一个新流程"
  >
    <UButton
      label="创建流程"
      icon="i-lucide-git-pull-request-create-arrow"
      color="neutral"
    />
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
            autofocus
          />
        </UFormField>
        <!-- 细节描述 -->
        <UFormField label="描述" name="description">
          <UInput
            v-model="state.description"
            class="w-full"
            placeholder="请输入细节描述"
          />
        </UFormField>
        <!-- 负责人+流程序号 -->
        <div class="grid grid-cols-2 gap-4">
          <UFormField label="负责人" name="employee" required>
            <USelectMenu
              v-model="state.employee"
              :items="employeeItems"
              value-key="value"
              placeholder="请选择负责人"
              class="w-full"
            />
          </UFormField>
          <UFormField label="流程序号" name="sequence">
            <UInputNumber
              v-model="state.sequence"
              class="w-full"
              placeholder="请输入流程序号"
              :min="1"
              :max="maxSequence"
            />
          </UFormField>
        </div>
        <!-- 备注 -->
        <UFormField label="备注" name="remark">
          <UTextarea
            v-model="state.remark"
            autoresize
            class="w-full"
            placeholder="请输入备注"
          />
        </UFormField>
        <!-- 按钮 -->
        <div class="flex justify-end gap-2">
          <UButton
            label="取消"
            variant="outline"
            @click="open = false"
            @pointerdown.prevent
          />
          <UButton
            label="创建"
            type="submit"
            :loading="isCreating"
            color="neutral"
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
import type { InsertProcessRequest } from "../../models/processes"
import { storeToRefs } from "pinia"
import { getUserID } from "../../utils/auth"
import { useEmployeeStore } from "../../stores/employee"
import { useProcesses } from "../../composables/useProcesses"
import { useProcessStore } from "../../stores/process"

const open = ref(false)
const toast = useToast()

const employeeStore = useEmployeeStore()
const processStore = useProcessStore()
const { employeeItems } = storeToRefs(employeeStore)
const { isCreating, createProcess } = useProcesses()
const { maxSequence } = storeToRefs(processStore)

// 定义表单数据
const schema = z.object({
  name: z.string().nonempty("请输入流程名称"),
  sequence: z.number().optional(),
  employee: z.number(),
  description: z.string(),
  remark: z.string()
})

type Schema = z.output<typeof schema>
const userID = getUserID()

// 表单默认数据状态
const initialState: Schema = {
  name: "",
  sequence: undefined,
  employee: userID!,
  description: "",
  remark: ""
}
const state = reactive<Schema>({ ...initialState })

// 创建流程
async function handleSubmit(event: FormSubmitEvent<Schema>) {
  const req: Omit<InsertProcessRequest, "projectID"> = {
    name: event.data.name,
    employeeID: event.data.employee,
    description: event.data.description,
    remark: event.data.remark
  }
  try {
    await createProcess({ sequence: event.data.sequence, req })
    toast.add({
      title: `流程 ${event.data.name} 创建成功`,
      color: "success"
    })
    open.value = false
  } catch (err) {
    toast.add({
      title: "创建失败",
      color: "error"
    })
    console.log("创建流程失败:", err)
  }
}

// 监视弹窗，开始时重置表单
watch(open, (isOpen) => {
  if (isOpen) {
    Object.assign(state, initialState)
  }
})
</script>
