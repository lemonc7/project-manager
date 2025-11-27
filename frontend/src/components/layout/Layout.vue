<template>
  <SidebarProvider>
    <!-- 菜单栏 -->
    <Aside />
    <!-- 外侧容器 -->
    <div class="flex flex-1 flex-col bg-gray-50 p-4">
      <!-- 内侧圆角容器 -->
      <div
        class="flex flex-1 flex-col rounded-xl bg-white shadow overflow-hidden px-4"
      >
        <Header />
        <RouterView />
      </div>
    </div>
    <!-- toast容器 -->
  </SidebarProvider>
</template>

<script setup lang="ts">
</script>
