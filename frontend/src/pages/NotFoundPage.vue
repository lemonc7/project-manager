<template>
  <div
    class="flex flex-1 flex-col items-center justify-center text-center px-4"
  >
    <h1 class="text-6xl font-bold text-gray-800 mb-4 select-none">404</h1>
    <p class="text-lg text-gray-600 mb-6">抱歉，你访问的页面不存在。</p>
    <UButton color="neutral" @click="router.replace('/')"> 返回首页 </UButton>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router"

const router = useRouter()
</script>
