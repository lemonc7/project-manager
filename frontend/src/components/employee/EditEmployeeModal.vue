<template>
  <UModal
    :open="props.open"
    @update:open="onOpenChange"
    title="编辑"
    description="更新成员信息(不包括用户名和密码)"
  >
    <template #body>
      <UForm
        :schema="schema"
        :state="state"
        class="space-y-4"
        @submit="handleSubmit"
      >
        <!-- 姓名 -->
        <UFormField label="姓名" name="name" required>
          <UInput
            v-model="state.name"
            class="w-full"
            placeholder="请输入姓名"
          />
        </UFormField>
        <!-- 小组 -->
        <UFormField label="小组" name="team">
          <USelect
            v-model="state.team"
            :items="teamItems"
            placeholder="请选择小组"
            class="w-full"
          />
        </UFormField>
        <!-- 角色 -->
        <UFormField label="角色" name="roles">
          <USelect
            v-model="state.roles"
            :items="[...roleItems]"
            placeholder="请选择角色"
            multiple
            class="w-full"
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
          <UButton
            label="更新"
            type="submit"
            :loading="isUpdating"
            color="neutral"
          />
        </div>
      </UForm>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import z from "zod"
import {
  roleItems,
  type Employee,
  type UpdateEmployeeRequest
} from "../../models/employees"
import { reactive, watch } from "vue"
import { storeToRefs } from "pinia"
import type { FormSubmitEvent } from "@nuxt/ui"
import { useTeamStore } from "../../stores/team"
import { useEmployees } from "../../composables/useEmployees"

const props = defineProps<{
  open: boolean
  employee: Employee | null
}>()

const emit = defineEmits<{
  (e: "update:open", open: boolean): void
}>()

const teamStore = useTeamStore()
const { teamItems } = storeToRefs(teamStore)
const { updateEmployee, isUpdating } = useEmployees()

const toast = useToast()

// 定义表单数据
const schema = z.object({
  name: z.string().nonempty("请输入姓名"),
  roles: z.array(z.enum(["admin", "manager", "employee"])).default([]),
  team: z.number().optional()
})

type Schema = z.output<typeof schema>

const state = reactive<Schema>({
  name: "",
  roles: [],
  team: undefined
})

function onOpenChange(open: boolean) {
  emit("update:open", open)
}

async function handleSubmit(event: FormSubmitEvent<Schema>) {
  if (!props.employee) return
  try {
    const req: UpdateEmployeeRequest = {
      name: event.data.name,
      roles: event.data.roles,
      teamID: event.data.team
    }

    await updateEmployee({ id: props.employee.id, req })
    toast.add({
      title: `成员 ${event.data.name} 更新成功`,
      color: "success"
    })
    onOpenChange(false)
  } catch (err) {
    toast.add({
      title: "更新失败",
      color: "error"
    })
    console.log("更新成员失败:", err)
  }
}

// 监听open变化，打开时同步表单数据
watch(
  () => props.open,
  (isOpen) => {
    if (isOpen && props.employee) {
      Object.assign(state, {
        name: props.employee.name,
        team: props.employee.teamID,
        roles: props.employee.roles
      })
    }
  }
)
</script>
