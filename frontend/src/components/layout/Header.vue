<template>
  <header
    class="sticky top-0 z-10 flex h-16 items-center justify-between bg-background"
  >
    <!-- 左侧aside折叠按钮|标题|主页|上一页 下一页 -->
    <div class="flex items-center space-x-2 h-4">
      <SidebarTrigger class="w-10 h-10" title="打开/关闭菜单" />
      <Separator orientation="vertical" class="bg-gray-400" />
      <h1 class="font-semibold">{{ title }}</h1>
      <Separator orientation="vertical" class="bg-gray-400" />
      <Button
        title="主页"
        variant="outline"
        size="sm"
        class="ml-2"
        @click="router.push('/home')"
      >
        <Home />
      </Button>
      <ButtonGroup class="hidden md:flex">
        <Button
          title="后退"
          variant="outline"
          size="sm"
          class="w-12"
          @click="router.go(-1)"
        >
          <ArrowLeft />
        </Button>
        <Button
          title="前进"
          variant="outline"
          size="sm"
          class="w-12"
          @click="router.go(1)"
        >
          <ArrowRight />
        </Button>
      </ButtonGroup>
    </div>
    <!-- 右侧动态元素|刷新 -->
    <div class="flex items-center gap-2">
      <div id="header-actions" class="flex items-center gap-2"></div>
      <Button variant="outline" size="sm" title="刷新" @click="handleRefresh">
        <RefreshCw />
      </Button>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useQueryClient } from "@tanstack/vue-query"
import { ArrowLeft, ArrowRight, Home, RefreshCw } from "lucide-vue-next"
import { useRoute, useRouter } from "vue-router"
import { toast } from "vue-sonner"

const route = useRoute()
const router = useRouter()
const queryClient = useQueryClient()

async function handleRefresh() {
  await queryClient.refetchQueries({ type: "active" })
  toast.success("页面已刷新")
}

// 从路由meta里取标题
const title = route.meta?.title ?? ""
</script>
