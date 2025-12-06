<template>
  <UModal
    v-model:open="open"
    title="创建"
    description="添加一个新成员到数据库中"
  >
    <UButton label="创建成员" icon="i-lucide-user-round-plus" color="neutral" />
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
        <!-- 小组+角色 -->
        <div class="grid grid-cols-2 gap-4">
          <UFormField label="小组" name="team">
            <USelect
              v-model="state.team"
              :items="teamItems"
              item-label="label"
              item-value="value"
              placeholder="请选择小组"
              class="w-full"
            />
          </UFormField>
          <UFormField label="角色" name="roles">
            <USelect
              v-model="state.roles"
              :items="[...roleItems]"
              item-label="label"
              item-value="value"
              placeholder="请选择角色"
              multiple
              class="w-full"
            />
          </UFormField>
        </div>
        <!-- 用户名 -->
        <UFormField label="用户名" name="username" required>
          <UInput
            v-model="state.username"
            class="w-full"
            placeholder="请输入用户名"
          />
        </UFormField>
        <!-- 密码 -->
        <UFormField label="密码" name="password" required>
          <UInput
            v-model="state.password"
            :type="show ? 'text' : 'password'"
            class="w-full"
            placeholder="请输入密码"
            :ui="{ trailing: 'pe-1' }"
          >
            <template #trailing>
              <UButton
                color="neutral"
                variant="link"
                size="sm"
                :icon="show ? 'i-lucide-eye-off' : 'i-lucide-eye'"
                :aria-label="show ? 'Hide password' : 'Show password'"
                :aria-pressed="show"
                aria-controls="password"
                @click="show = !show"
              />
            </template>
          </UInput>
        </UFormField>
        <!-- 按钮 -->
        <div class="flex justify-end gap-2">
          <UButton
            label="取消"
            color="neutral"
            variant="outline"
            @click="open = false"
          />
          <UButton
            label="创建"
            type="submit"
            color="neutral"
            :loading="createEmployee.isPending.value"
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
import useCreateEmployee from "../../composables/useCreateEmployee"
import { roleItems, type CreateEmployeeRequest } from "../../models/employees"

// 需要传进来team的选项
defineProps<{
  teamItems: {
    label: string
    value: number
  }[]
}>()

const open = ref(false)
const show = ref(false)

// 定义表单数据
const schema = z.object({
  name: z.string().nonempty("请输入姓名"),
  username: z.string().nonempty("请输入用户名"),
  password: z.string().min(6, "密码不能少于6位"),
  team: z.number().optional(),
  roles: z.array(z.enum(["admin", "manager", "employee"] as const))
})

type Schema = z.output<typeof schema>

// 表单默认数据状态
const initialState: Schema = {
  name: "",
  username: "",
  password: "",
  team: undefined,
  roles: ["employee"]
}
const state = reactive<Schema>({ ...initialState })

const toast = useToast()
// 创建员工
const createEmployee = useCreateEmployee({
  onSuccess: () => {
    toast.add({
      title: `成员 ${state.name} 创建成功`,
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
  const req: CreateEmployeeRequest = {
    name: event.data.name,
    username: event.data.username,
    password: event.data.password,
    teamID: event.data.team,
    roles: event.data.roles
  }

  createEmployee.mutate(req)
}

// 监视弹窗，关闭时重置表单
watch(open, (newVal) => {
  if (!newVal) {
    Object.assign(state, initialState)
  }
})
</script>
