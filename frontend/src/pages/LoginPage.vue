<template>
  <div class="flex min-h-screen items-center justify-center bg-gray-100">
    <UPageCard class="w-full max-w-md">
      <UAuthForm
        :schema="schema"
        :fields="fields"
        title="登录"
        icon="i-lucide-user"
        @submit="handleSubmit"
        :submit="{
          label: '登录',
        }"
      >
        <template #validation v-if="error !== ''">
          <UAlert color="error" icon="i-lucide-info" :title="error" />
        </template>
      </UAuthForm>
    </UPageCard>
  </div>
</template>

<script setup lang="ts">
import type { AuthFormField, FormSubmitEvent } from "@nuxt/ui"
import { login } from "../api/auth"
import { ref } from "vue"
import z from "zod"
import { useRouter } from "vue-router"

const router = useRouter()
const error = ref("")

const fields = ref<AuthFormField[]>([
  {
    name: "username",
    type: "text",
    label: "用户名",
    placeholder: "请输入用户名",
    required: true
  },
  {
    name: "password",
    type: "password",
    label: "密码",
    placeholder: "请输入密码",
    required: true
  }
])

const schema = z.object({
  username: z.string("需要输入用户名"),
  password: z.string("需要输入密码")
})

type Schema = z.output<typeof schema>

const handleSubmit = async (auth: FormSubmitEvent<Schema>) => {
  try {
    await login({ username: auth.data.username, password: auth.data.password })
    router.replace("/")
  } catch (e) {
    if (e instanceof Error) {
      error.value = e.message
    } else {
      error.value = String(e) || "未知错误"
    }
  }
}
</script>
