<template></template>

<script setup lang="ts">
import type { Employee, Role } from "@/models/employees"
import type { Team } from "@/models/teams"
import { reactive } from "vue"

const props = defineProps<{
  open: boolean
  teams: Team[]
  employee: Employee | null
}>()

const emit = defineEmits<{
  (e: "update:open", value: boolean): void
}>()

// 默认值
const formData = reactive<{
  name: string
  role: Role
  teamID?: number
}>({
  name: "",
  role: "employee"
})
</script>
