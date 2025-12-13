<template>
  <UModal
    :open="open"
    @update:open="onOpenChange"
    title="复制"
    description="复制项目与相关流程，并重置状态"
  >
    <template #body>
      <UForm
        :schema="schema"
        :state="state"
        class="space-y-4"
        @submit="handleSubmit"
      >
        <!-- 项目名称 -->
        <UFormField label="新名称" name="name" required>
          <UFieldGroup class="w-full">
            <UInput
              v-model="state.name"
              class="w-full"
              placeholder="请输入新的项目名称"
            />
            <UButton
              icon="i-lucide-check"
              variant="outline"
              label="检查"
              @click="handleCheck"
              :loading="isChecking"
            />
          </UFieldGroup>
        </UFormField>
        <!-- 按钮 -->
        <div class="flex justify-end gap-2">
          <UButton
            label="取消"
            variant="outline"
            @click="onOpenChange(false)"
            @pointerdown.prevent
          />
          <UButton
            label="确定"
            type="submit"
            :loading="isCopying"
            color="neutral"
          />
        </div>
      </UForm>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import { storeToRefs } from "pinia"
import { useProjectStore } from "../../composables/useProjectStore"
import z from "zod"
import { getUserID } from "../../utils/auth"
import { reactive } from "vue"
import type { FormSubmitEvent } from "@nuxt/ui"
import type { CopyProjectRequest, Project } from "../../models/projects"

const props = defineProps<{
  open: boolean
  project: Project | null
}>()

const emit = defineEmits<{
  (e: "update:open", open: boolean): void
}>()

function onOpenChange(open: boolean) {
  emit("update:open", open)
}

const toast = useToast()

// 定义表单数据
const schema = z.object({
  name: z.string().nonempty("请输入新的项目名称")
})

type Schema = z.output<typeof schema>

const userID = getUserID()
const initialState: Schema = {
  name: ""
}

const state = reactive<Schema>({ ...initialState })

const projectStore = useProjectStore()
const { copyProject, checkProject } = projectStore
const { isCopying, isChecking } = storeToRefs(projectStore)

// 表单提交
async function handleSubmit(event: FormSubmitEvent<Schema>) {
  if (!props.project) return
  const req: CopyProjectRequest = {
    newName: event.data.name,
    managerID: userID!
  }

  try {
    await copyProject({ id: props.project.id, req })
    toast.add({
      title: `项目 ${event.data.name} 复制成功`,
      color: "success"
    })
    onOpenChange(false)
  } catch (err) {
    toast.add({
      title: "复制失败",
      color: "error"
    })
    console.log("复制项目失败:", err)
  }
}

// 检查项目名称是否存在
async function handleCheck() {
  try {
    const exist = await checkProject(state.name)
    if (exist) {
      toast.add({
        title: `项目 ${state.name} 已存在`,
        color: "error"
      })
    } else {
      toast.add({
        title: "项目名称可用",
        color: "success"
      })
    }
  } catch (err) {
    toast.add({
      title: "检查失败",
      color: "error"
    })
    console.log("检查项目名称失败:", err)
  }
}
</script>
