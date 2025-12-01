<template>
  <Dialog :open="props.open" @update:open="onOpenChange">
    <DialogContent>
      <DialogTitle>编辑成员</DialogTitle>
      <DialogDescription></DialogDescription>
      <div class="w-full max-w-md">
        <form @submit.prevent="handleSubmit">
          <FieldGroup>
            <FieldSet>
              <!-- 姓名 -->
              <Field>
                <FieldLabel for="employee-name">姓名</FieldLabel>
                <Input
                  id="employee-name"
                  v-model="formData.name"
                  placeholder="请输入成员名称"
                  required
                />
              </Field>
              <!-- 小组 -->
              <Field>
                <FieldLabel>小组</FieldLabel>
                <Select v-model="formData.teamID">
                  <SelectTrigger>
                    <SelectValue placeholder="选择小组" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectGroup>
                      <SelectItem
                        v-for="team in teams"
                        :key="team.id"
                        :value="team.id"
                      >
                        {{ team.name }}
                      </SelectItem>
                    </SelectGroup>
                  </SelectContent>
                </Select>
              </Field>
              <!-- 角色 -->
              <Field>
                <FieldLabel>角色</FieldLabel>
                <div class="flex flex-wrap gap-4">
                  <div
                    class="flex items-center space-x-2 mt-2"
                    v-for="r in roles"
                    :key="r"
                  >
                    <input
                      type="checkbox"
                      v-model="formData.roles"
                      :value="r"
                      :id="r"
                    />
                    <Label :for="r">{{ RoleMap[r] }}</Label>
                  </div>
                </div>
              </Field>
            </FieldSet>
            <Field orientation="horizontal" class="justify-end">
              <Button
                type="button"
                variant="outline"
                @click="onOpenChange(false)"
                >取消</Button
              >
              <Button
                type="submit"
                :disabled="updateEmployee.isPending.value"
                >{{
                  updateEmployee.isPending.value ? "更新中..." : "更新"
                }}</Button
              >
            </Field>
          </FieldGroup>
        </form>
      </div>
    </DialogContent>
  </Dialog>
</template>

<script setup lang="ts">
import useUpdateEmployee from "@/composables/use-update-employee"
import {
  RoleMap,
  type Employee,
  type Role,
  type UpdateEmployeeRequest
} from "@/models/employees"
import type { Team } from "@/models/teams"
import { reactive, watch } from "vue"
import { toast } from "vue-sonner"

const props = defineProps<{
  open: boolean
  teams: Team[]
  employee: Employee | null
}>()

const emit = defineEmits<{
  (e: "update:open", value: boolean): void
}>()

const roles = Object.keys(RoleMap) as Role[]

// 默认值
const formData = reactive<{
  name: string
  roles?: Role[]
  teamID?: number
}>({
  name: "",
  roles: ["employee"]
})

const onOpenChange = (value: boolean) => {
  emit("update:open", value)
}

const updateEmployee = useUpdateEmployee({
  onSuccess: () => {
    toast.success(`成员 ${formData.name} 更新成功`)
    onOpenChange(false)
  },
  onError: () => {
    toast.error("更新失败")
  }
})

// 监听props.employee变化，用于数据初始化和外部数据更新
watch(
  () => props.employee,
  (newEmployee) => {
    if (newEmployee) {
      Object.assign(formData, {
        name: newEmployee.name,
        roles: newEmployee.roles,
        teamID: newEmployee.teamID
      })
    }
  },
  { immediate: true }
)

// 监听open变化，用于取消操作后的数据重置
watch(
  () => props.open,
  (isOpen) => {
    if (!isOpen && props.employee) {
      // 弹窗关闭时，将formData重置为原始的employee值
      Object.assign(formData, {
        name: props.employee.name,
        roles: props.employee.roles,
        teamID: props.employee.teamID
      })
    }
  }
)

// 表单提交处理
const handleSubmit = () => {
  if (!props.employee?.id) return
  // 基础校验
  if (!formData.name) {
    toast.error("请填写成员名称")
    return
  }

  // 构造请求
  const req: UpdateEmployeeRequest = {
    name: formData.name,
    roles: formData.roles,
    teamID: formData.teamID
  }

  updateEmployee.mutate({ id: props.employee.id, req })
}
</script>
