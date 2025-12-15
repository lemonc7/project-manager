<template>
  <UModal
    :open="props.open"
    @update:open="onOpenChange"
    title="重置密码"
    description="重置用户的密码"
  >
    <template #body>
      <UForm
        :schema="schema"
        :state="state"
        class="space-y-4"
        @submit="handleSubmit"
      >
        <!-- 密码 -->
        <UFormField label="密码" name="password" required>
          <InputPassword v-model="state.password" placeholder="请输入密码" />
        </UFormField>
        <div class="flex justify-end gap-2">
          <UButton
            label="取消"
            color="neutral"
            @click="onOpenChange(false)"
          />
          <UButton label="确认" variant="neutral" type="submit" />
        </div>
      </UForm>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import type { FormSubmitEvent } from "@nuxt/ui"
import { reactive, watch } from "vue"
import z from "zod"
import { resetPassword } from "../../api/auth"
import type { Employee } from "../../models/employees"

const props = defineProps<{
  open: boolean
  employee: Employee | null
}>()

const emit = defineEmits<{
  (e: "update:open", open: boolean): void
}>()

const toast = useToast()

const schema = z.object({
  password: z.string().min(6, "密码不能少于6位").max(20, "密码不能超过20位")
})
type Schema = z.output<typeof schema>

const initialValues: Schema = {
  password: ""
}
const state = reactive<Schema>({ ...initialValues })

function onOpenChange(open: boolean) {
  emit("update:open", open)
}

async function handleSubmit(event: FormSubmitEvent<Schema>) {
  if (!props.employee) return

  try {
    await resetPassword(props.employee.id, event.data.password)
    toast.add({
      title: "密码重置成功",
      color: "success"
    })
    onOpenChange(false)
  } catch (err) {
    toast.add({
      title: "重置失败",
      color: "error"
    })
    console.log("重置密码失败:", err)
  }
}

watch(
  () => props.open,
  (isOpen) => {
    if (isOpen) {
      Object.assign(state, initialValues)
    }
  }
)
</script>
