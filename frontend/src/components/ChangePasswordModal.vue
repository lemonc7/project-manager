<template>
  <UModal
    :open="props.open"
    @update:open="onOpenChange"
    title="修改密码"
    description="修改用户密码，密码位数为6~20"
  >
    <template #body>
      <UForm
        :schema="schema"
        :state="state"
        class="space-y-4"
        @submit="handleSubmit"
      >
        <!-- 旧密码 -->
        <UFormField label="原密码" name="oldPassword" required>
          <InputPassword v-model="state.oldPassword" placeholder="请输入密码" />
        </UFormField>
        <!-- 新密码 -->
        <UFormField label="新密码" name="newPassword" required>
          <InputPassword v-model="state.newPassword" placeholder="请输入密码" />
        </UFormField>
        <!-- 确认密码 -->
        <UFormField label="确认密码" name="confirmPassword" required>
          <InputPassword
            v-model="state.confirmPassword"
            placeholder="再次输入密码"
          />
        </UFormField>
        <!-- 按钮 -->
        <div class="flex justify-end gap-2">
          <UButton
            label="取消"
            variant="outline"
            @click="onOpenChange(false)"
          />
          <UButton label="确认" type="submit" color="neutral" />
        </div>
      </UForm>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import type { FormSubmitEvent } from "@nuxt/ui"
import { reactive, watch } from "vue"
import z from "zod"
import { changePassword } from "../api/auth"
import { removeToken } from "../utils/auth"

const props = defineProps<{
  open: boolean
}>()

const emit = defineEmits<{
  (e: "update:open", open: boolean): void
}>()

const toast = useToast()

// 定义表单数据
const schema = z.object({
  oldPassword: z.string().nonempty("请输入旧密码"),
  newPassword: z.string().min(6, "密码不能少于6位").max(20, "密码不能超过20位"),
  confirmPassword: z
    .string()
    .min(6, "密码不能少于6位")
    .max(20, "密码不能超过20位")
})

type Schema = z.output<typeof schema>

// 表单默认数据状态
const initialState: Schema = {
  oldPassword: "",
  newPassword: "",
  confirmPassword: ""
}
const state = reactive<Schema>({ ...initialState })

function onOpenChange(open: boolean) {
  emit("update:open", open)
}

// 表单提交
async function handleSubmit(event: FormSubmitEvent<Schema>) {
  if (event.data.newPassword !== event.data.confirmPassword) {
    toast.add({
      title: "输入的两次新密码不一致",
      color: "error"
    })
    return
  }

  try {
    await changePassword({
      oldPassword: event.data.oldPassword,
      newPassword: event.data.newPassword
    })
    toast.add({
      title: "修改成功,请重新登录",
      color: "success"
    })
    onOpenChange(false)
    removeToken()
  } catch (err) {
    toast.add({
      title: "修改失败",
      color: "error"
    })
    console.log("修改密码失败:", err)
  }
}

// 监视弹窗，打开是重置表单
watch(
  () => props.open,
  (isOpen) => {
    if (isOpen) {
      Object.assign(state, initialState)
    }
  }
)
</script>
