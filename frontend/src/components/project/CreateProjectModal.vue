<template>
  <UModal
    v-model:open="open"
    title="创建"
    description="添加一个新项目到数据库中"
  >
    <UButton
      label="创建项目"
      icon="i-lucide-calendar-plus"
      color="neutral"
      :disabled="!hasAtLeastRole('manager')"
    />
    <template #body>
      <UForm
        :schema="schema"
        :state="state"
        class="space-y-4"
        color="neutral"
        @submit="handleSubmit"
      >
        <!-- 项目名称 -->
        <UFormField label="名称" name="name" required>
          <UInput
            v-model="state.name"
            class="w-full"
            placeholder="请输入项目名称"
            autofocus
          />
        </UFormField>
        <UFormField label="描述" name="description">
          <UInput
            v-model="state.description"
            class="w-full"
            placeholder="请输入细节描述"
          />
        </UFormField>
        <!-- 优先级+项目周期 -->
        <div class="grid grid-cols-2 gap-4">
          <UFormField label="优先级" name="priority" required>
            <UInputNumber
              v-model="state.priority"
              class="w-full"
              placeholder="请输入优先级"
              :min="1"
              :max="5"
            />
          </UFormField>
          <UFormField label="项目周期" name="duration">
            <UInputNumber
              v-model="state.duration"
              class="w-full"
              placeholder="请输入项目周期"
              :min="0"
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
import type { CreateProjectRequest } from "../../models/projects"
import { getUserID, hasAtLeastRole } from "../../utils/auth"
import { useProjectStore } from "../../composables/useProjectStore"
import { storeToRefs } from "pinia"

const open = ref(false)

// 定义表单数据
const schema = z.object({
  name: z.string().nonempty("请输入项目名称"),
  description: z.string().default(""),
  remark: z.string().default(""),
  priority: z.number().default(1),
  duration: z.number().min(0, "项目周期不能小于0").optional()
})

type Schema = z.output<typeof schema>

const userID = getUserID()
// 表单默认数据状态
const initialState: Schema = {
  name: "",
  description: "",
  remark: "",
  priority: 1,
  duration: undefined
}
const state = reactive<Schema>({ ...initialState })

const toast = useToast()

const projectStore = useProjectStore()
const { createProject } = projectStore
const { isCreating } = storeToRefs(projectStore)

// 表单提交
async function handleSubmit(event: FormSubmitEvent<Schema>) {
  const req: CreateProjectRequest = {
    name: event.data.name,
    description: event.data.description,
    managerID: userID!,
    remark: event.data.remark,
    priority: event.data.priority,
    duration: event.data.duration
  }

  try {
    await createProject(req)
    toast.add({
      title: `项目 ${event.data.name} 创建成功`,
      color: "success"
    })
    open.value = false
  } catch (err) {
    toast.add({
      title: "创建失败",
      color: "error"
    })
    console.log("创建项目失败:", err)
  }
}

// 监视弹窗，打开时重置表单
watch(open, (isOpen) => {
  if (isOpen) {
    Object.assign(state, initialState)
  }
})
</script>
