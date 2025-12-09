<template>
  <UModal
    :open="props.open"
    @update:open="onOpenChange"
    title="编辑"
    description="更新小组信息"
  >
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
          />
        </UFormField>
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
            @click="onOpenChange(false)"
            @pointerdown.prevent
          />
          <UButton label="更新" type="submit" :loading="isUpdating" />
        </div>
      </UForm>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import z from "zod"
import { reactive, watch } from "vue"
import { useTeamStore } from "../../composables/useTeamStore"
import { storeToRefs } from "pinia"
import type { Team, UpdateTeamRequest } from "../../models/teams"
import type { FormSubmitEvent } from "@nuxt/ui"

const props = defineProps<{
  open: boolean
  team: Team | null
}>()

const emit = defineEmits<{
  (e: "update:open", open: boolean): void
}>()

const teamStore = useTeamStore()
const { isUpdating } = storeToRefs(teamStore)
const { updateTeam } = teamStore

const toast = useToast()

// 定义表单数据
const schema = z.object({
  name: z.string().nonempty("请输入名称"),
  description: z.string().default("")
})

type Schema = z.output<typeof schema>

const state = reactive<Schema>({
  name: "",
  description: ""
})

function onOpenChange(open: boolean) {
  emit("update:open", open)
}

async function handleSubmit(event: FormSubmitEvent<Schema>) {
  if (!props.team) return
  try {
    const req: UpdateTeamRequest = {
      name: event.data.name,
      description: event.data.description
    }

    await updateTeam({ id: props.team.id, req })
    toast.add({
      title: `小组 ${event.data.name} 更新成功`,
      color: "success"
    })
    onOpenChange(false)
  } catch (err) {
    toast.add({
      title: "更新失败",
      color: "error"
    })
    console.log("更新小组失败:", err)
  }
}

// 监听open变化，打开时同步表单数据
watch(
  () => props.open,
  (isOpen) => {
    if (isOpen && props.team) {
      Object.assign(state, {
        name: props.team.name,
        description: props.team.description
      })
    }
  }
)
</script>
