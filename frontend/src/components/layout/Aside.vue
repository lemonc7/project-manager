<template>
  <div>
    <Sidebar class="border-none!">
      <!-- Header -->
      <SidebarHeader class="mt-2">
        <Button variant="ghost" @click="router.replace('/home')">
          <User />
          项目管理平台
        </Button>
      </SidebarHeader>
      <!-- Content -->
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>项目管理</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem v-for="item in items" :key="item.title">
                <SidebarMenuButton as-child>
                  <RouterLink :to="item.url">
                    <component :is="item.icon" />
                    <span>{{ item.title }}</span>
                  </RouterLink>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <!-- Footer -->
      <SidebarFooter class="flex justify-start py-4">
        <DropdownMenu>
          <DropdownMenuTrigger as-child>
            <Button variant="ghost">
              <User />
              <span>{{ data?.name ?? "未登录" }}</span>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuItem @click="logout">
              <LogOut />
              登出
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarFooter>
    </Sidebar>
  </div>
</template>

<script setup lang="ts">
import { getEmployeeMe } from "@/api/employees"
import { removeToken } from "@/util/auth"
import { useQuery } from "@tanstack/vue-query"
import { RouterLink, useRouter } from "vue-router"
import { ListTodo, LogOut, Sheet, User } from "lucide-vue-next"

const router = useRouter()

const items = [
  { title: "任务", url: "/home", icon: ListTodo },
  { title: "项目总览", url: "/projects", icon: Sheet },
  { title: "成员", url: "/employees", icon: User }
]

const { data } = useQuery({
  queryKey: ["user"],
  queryFn: getEmployeeMe
})

function logout() {
  removeToken()
  router.replace("/login")
}
</script>
