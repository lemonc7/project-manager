<template>
  <UModal
    :open="props.open"
    @update:open="onOpenChange"
    title="编辑成员"
    description="更新成员的内容(不包括用户名和密码)"
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
          />
          <UButton label="更新" type="submit" :loading="isUpdating" />
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
import { useEmployeeStore } from "../../composables/useEmployeeStore"
import { useTeamStore } from "../../composables/useTeamStore"
import { storeToRefs } from "pinia"

const props = defineProps<{
  open: boolean
  employee: Employee | null
}>()

const emit = defineEmits<{
  (e: "update:open", open: boolean): void
}>()

const teamStore = useTeamStore()
const employeeStore = useEmployeeStore()
const { teamItems } = storeToRefs(teamStore)
const { isUpdating } = storeToRefs(employeeStore)
const { updateEmployee } = employeeStore

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
  roles: []
})

function onOpenChange(open: boolean) {
  emit("update:open", open)
}

async function handleSubmit() {
  if (!props.employee?.id) return
  try {
    const req: UpdateEmployeeRequest = {
      name: state.name,
      roles: state.roles,
      teamID: state.team
    }

    await updateEmployee({ id: props.employee.id, req })
    toast.add({
      title: `成员 ${state.name} 更新成功`,
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

// 监听employee变化，用于数据初始化和外部数据更新
watch(
  () => props.employee,
  (newVal) => {
    if (newVal) {
      Object.assign(state, {
        name: newVal.name,
        team: newVal.teamID,
        roles: newVal.roles
      })
    }
  },
  { immediate: true }
)

// 监听open变化，关闭时重置表单
watch(
  () => props.open,
  (newVal) => {
    if (!newVal && props.employee) {
      Object.assign(state, {
        name: props.employee.name,
        team: props.employee.teamID,
        roles: props.employee.roles
      })
    }
  }
)
</script>
