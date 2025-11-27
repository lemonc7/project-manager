<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100">
    <Card class="w-full max-w-sm">
      <CardHeader>
        <CardTitle>登录</CardTitle>
      </CardHeader>
      <form @submit.prevent="handleSubmit">
        <CardContent>
          <div class="flex flex-col gap-6">
            <div class="grid gap-2">
              <Label for="username">用户名</Label>
              <Input
                id="username"
                v-model="username"
                placeholder="请输入用户名"
                required
              />
            </div>
            <div class="grid gap-2">
              <Label for="password">密码</Label>
              <Input
                id="password"
                type="password"
                v-model="password"
                placeholder="请输入密码"
                required
              />
            </div>
            <p v-if="err" class="text-red-500">{{ err }}</p>
          </div>
        </CardContent>
        <CardFooter class="flex-col gap-2 mt-6">
          <Button type="submit" class="w-full" :disabled="loading">
            {{ loading ? "登录中..." : "登录" }}
          </Button>
        </CardFooter>
      </form>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { login } from "@/api/auth"
import { ref } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()

const username = ref("")
const password = ref("")
const loading = ref(false)
const err = ref("")

async function handleSubmit() {
  loading.value = true

  try {
    await login({ username: username.value, password: password.value })
    router.replace("/home")
  } catch (e) {
    if (e instanceof Error) {
      err.value = e.message
    } else {
      err.value = String(e) || "未知错误"
    }
  } finally {
    loading.value = false
  }
}
</script>
