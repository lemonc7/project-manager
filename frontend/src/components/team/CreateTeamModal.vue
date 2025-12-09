<template>
  <UModal
    v-model:open="open"
    title="创建"
    description="添加一个新小组到数据库中"
  >
    <UButton label="创建小组" icon="fluent:people-team-add-24-regular" />
    <template #body>
      <UForm
        :schema="schema"
        :state="state"
        class="space-y-4"
        @submit="handleSubmit"
      >
        <!-- 名称 -->
        <UFormField label="名称" name="name" required>
          <UInput
            v-model="state.name"
            class="w-full"
            placeholder="请输入名称"
            autofocus
          />
        </UFormField>
        <!-- 描述 -->
        <UFormField label="描述" name="description">
          <UInput
            v-model="state.description"
            class="w-full"
            placeholder="请输入描述"
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
          <UButton label="创建" type="submit" :loading="isCreating" />
        </div>
      </UForm>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import type { FormSubmitEvent } from "@nuxt/ui"
import { reactive, ref, watch } from "vue"
import z from "zod"
import { useTeamStore } from "../../composables/useTeamStore"
import { storeToRefs } from "pinia"
import type { CreateTeamRequest } from "../../models/teams"

const open = ref(false)
const toast = useToast()

const teamStore = useTeamStore()
const { isCreating } = storeToRefs(teamStore)
const { createTeam } = teamStore

// 定义表单数据
const schema = z.object({
  name: z.string().nonempty("请输入名称"),
  description: z.string().default("")
})

type Schema = z.output<typeof schema>

// 表单默认数据状态
const initialState: Schema = {
  name: "",
  description: ""
}
const state = reactive<Schema>({ ...initialState })

// 表单提交
async function handleSubmit(event: FormSubmitEvent<Schema>) {
  const req: CreateTeamRequest = {
    name: event.data.name,
    description: event.data.description
  }

  try {
    await createTeam(req)
    toast.add({
      title: `小组 ${event.data.name} 创建成功`,
      color: "success"
    })
    open.value = false
  } catch (err) {
    toast.add({
      title: "创建失败",
      color: "error"
    })
    console.log("创建小组失败:", err)
  }
}

// 监视弹窗，打开时重置表单
watch(open, (isOpen) => {
  if (isOpen) {
    Object.assign(state, initialState)
  }
})
</script>
