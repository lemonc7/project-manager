<template>
  <Dialog :open="props.open" @update:open="onOpenChange">
    <DialogContent>
      <DialogTitle>创建成员</DialogTitle>
      <DialogDescription></DialogDescription>
      <div class="w-full max-w-md">
        <form @submit.prevent="handleSubmit">
          <FieldGroup>
            <FieldSet>
              <div class="grid grid-cols-2 gap-4">
                <!-- 姓名 -->
                <Field>
                  <FieldLabel for="create-name">姓名</FieldLabel>
                  <Input
                    id="create-name"
                    v-model="formData.name"
                    placeholder="请输入成员姓名"
                    required
                  />
                </Field>
                <!-- 小组 -->
                <Field>
                  <FieldLabel for="create-team">小组</FieldLabel>
                  <Select v-model="formData.teamID">
                    <SelectTrigger>
                      <SelectValue placeholder="选择小组" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem
                        v-for="team in props.teams"
                        :key="team.id"
                        :value="team.id"
                      >
                        {{ team.name }}
                      </SelectItem>
                    </SelectContent>
                  </Select>
                </Field>
              </div>
              <!-- 用户名 -->
              <Field>
                <FieldLabel for="create-username">用户名</FieldLabel>
                <Input
                  id="create-username"
                  v-model="formData.username"
                  placeholder="请输入用户名"
                  required
                />
              </Field>
              <!-- 密码 -->
              <Field>
                <FieldLabel for="create-password">密码</FieldLabel>
                <Input
                  id="create-password"
                  v-model="formData.password"
                  placeholder="请输入密码"
                  required
                />
              </Field>
              <!-- 角色权限 -->
              <Field>
                <FieldLabel>角色权限</FieldLabel>
                <div class="flex flex-wrap gap-4">
                  <div
                    v-for="r in roles"
                    :key="r"
                    class="flex items-center space-x-2"
                  >
                    <Checkbox
                      :id="r"
                      :checked="isChecked(r)"
                      @update:checked="(checked: boolean) => toggleRole(r, checked)"
                    />
                    <Label :for="r">{{ RoleMap[r] }}</Label>
                  </div>
                </div>
              </Field>
            </FieldSet>
            <Field orientation="horizontal" class="justify-end">
              <Button
                variant="outline"
                type="button"
                @click="onOpenChange(false)"
                >取消</Button
              >
              <Button
                type="submit"
                :disabled="createEmployee.isPending.value"
                >{{
                  createEmployee.isPending.value ? "创建中..." : "创建"
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
import useCreateEmployee from "@/composables/use-create-employee"
import {
  RoleMap,
  type CreateEmployeeRequest,
  type Role
} from "@/models/employees"
import type { Team } from "@/models/teams"
import { reactive, watch } from "vue"
import { toast } from "vue-sonner"

const props = defineProps<{
  open: boolean
  teams: Team[]
}>()
const emit = defineEmits<{
  (e: "update:open", value: boolean): void
}>()

const roles = Object.keys(RoleMap) as Role[]

// 默认值
const initialFormData = {
  name: "",
  username: "",
  password: "",
  roles: ["employee"] as Role[],
  teamID: undefined
}

const formData = reactive({ ...initialFormData })

// 检查角色是否选中
const isChecked = (role: Role) => formData.roles.includes(role)

const toggleRole = (role: Role, checked: boolean) => {
  if (checked) {
    // 选中时，如果不在数组中就添加
    if (!formData.roles.includes(role)) {
      formData.roles.push(role)
    }
  } else {
    // 取消选中，从数组中移除
    formData.roles = formData.roles.filter((r) => r !== role)
  }
}

// 触发update:open事件，通知父组件关闭弹窗
const onOpenChange = (value: boolean) => {
  emit("update:open", value)
}

// 创建成员的方法
const createEmployee = useCreateEmployee({
  onSuccess: () => {
    toast.success(`成员 ${formData.name} 创建成功`)
    onOpenChange(false)
  },
  onError: () => {
    toast.error("创建失败")
  }
})

// 表单提交处理
const handleSubmit = () => {
  // 校验数据
  if (!formData.name || !formData.username || !formData.password) {
    toast.error("请填写姓名、用户名和密码")
    return
  }

  const req: CreateEmployeeRequest = {
    name: formData.name,
    username: formData.username,
    password: formData.password,
    roles: formData.roles,
    teamID: formData.teamID
  }

  createEmployee.mutate(req)
}

// 数据重置，监听open状态，在弹窗关闭时重置数据
watch(
  () => props.open,
  (newVal) => {
    if (!newVal) {
      Object.assign(formData, initialFormData)
    }
  }
)
</script>
