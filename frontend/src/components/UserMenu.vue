<template>
  <UDropdownMenu
    :items="items"
    :content="{ align: 'center', collisionPadding: 12 }"
    :ui="{
      content: collapsed ? 'w-48' : 'w-(--reka-dropdown-menu-trigger-width)'
    }"
  >
    <UButton
      v-bind:label="collapsed ? undefined : data?.name"
      v-bind:trailing-icon="collapsed ? undefined : 'i-lucide-chevrons-up-down'"
      variant="ghost"
      block
      icon="i-lucide-user"
      :square="collapsed"
      class="data-[state=open]:bg-elevated"
      :ui="{
        trailingIcon: 'text-dimmed'
      }"
    />
  </UDropdownMenu>
</template>

<script setup lang="ts">
import { useQuery } from "@tanstack/vue-query"
import { getEmployeeMe } from "../api/employees"
import type { DropdownMenuItem } from "@nuxt/ui"
import { computed } from "vue"
import { removeToken } from "../utils/auth"
import { useRouter } from "vue-router"

defineProps<{
  collapsed?: boolean
}>()

const { data } = useQuery({
  queryKey: ["user"],
  queryFn: getEmployeeMe
})

const router = useRouter()

const items = computed<DropdownMenuItem[]>(() => [
  {
    label: "更改密码",
    icon: "i-lucide-key"
  },
  {
    label: "登出",
    icon: "i-lucide-log-out",
    onSelect: () => {
      removeToken()
      router.replace("/login")
    }
  }
])
</script>
