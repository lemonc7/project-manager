<template>
  <Dialog :open="props.open" @update:open="onOpenChange">
    <DialogContent>
      <DialogTitle>创建项目</DialogTitle>
      <DialogDescription></DialogDescription>
      <div class="w-full max-w-md">
        <form @submit.prevent="handleSubmit">
          <FieldGroup>
            <FieldSet>
              <!-- 项目名称 -->
              <Field>
                <FieldLabel for="project-name">名称</FieldLabel>
                <Input
                  id="project-name"
                  v-model="formData.name"
                  placeholder="请输入项目名称"
                  required
                />
              </Field>
              <!-- 描述 -->
              <Field>
                <FieldLabel for="project-description">描述</FieldLabel>
                <Input
                  id="project-description"
                  v-model="formData.description"
                  placeholder="请输入细节描述"
                />
              </Field>
              <div class="grid grid-cols-3 gap-4">
                <!-- 负责人 -->
                <Field>
                  <FieldLabel>负责人</FieldLabel>
                  <EmployeesSelect v-model:employee-id="formData.managerID" />
                </Field>
                <!-- 优先级 -->
                <Field>
                  <FieldLabel for="project-priority">优先级</FieldLabel>
                  <Input
                    id="project-priority"
                    v-model="formData.priority"
                    placeholder="请输入项目优先级"
                    required
                    type="number"
                    min="1"
                    max="9"
                  />
                </Field>
                <!-- 项目周期 -->
                <Field>
                  <FieldLabel for="project-duration">周期</FieldLabel>
                  <Input
                    id="project-duration"
                    v-model="formData.duration"
                    placeholder="请输入项目周期"
                    type="number"
                    min="0"
                  />
                </Field>
              </div>
              <!-- 备注 -->
              <Field>
                <FieldLabel for="project-remark">备注</FieldLabel>
                <Textarea
                  id="project-remark"
                  placeholder="请输入项目备注"
                  class="resize-none"
                  v-model="formData.remark"
                />
              </Field>
            </FieldSet>
            <Field orientation="horizontal" class="justify-end">
              <Button
                variant="outline"
                type="button"
                @click="onOpenChange(false)"
                >取消</Button
              >
              <Button type="submit" :disabled="createProject.isPending.value">
                {{ createProject.isPending.value ? "创建中..." : "创建" }}
              </Button>
            </Field>
          </FieldGroup>
        </form>
      </div>
    </DialogContent>
  </Dialog>
</template>

<script setup lang="ts">
import EmployeesSelect from "@/components/common/EmployeesSelect.vue"
import useCreateProject from "@/composables/use-create-project"
import type { Employee } from "@/models/employees"
import type { CreateProjectRequest } from "@/models/projects"
import { getUserID } from "@/util/auth"
import { reactive, watch } from "vue"
import { toast } from "vue-sonner"

const props = defineProps<{
  open: boolean
  employees: Employee[]
}>()
const emit = defineEmits<{
  (e: "update:open", value: boolean): void
}>()

const userID = getUserID()

// 默认值
const initialFormData = {
  name: "",
  description: undefined,
  managerID: userID,
  remark: undefined,
  priority: 5,
  duration: undefined
}

const formData = reactive({ ...initialFormData })

// 触发update:open事件，通知父组件关闭弹窗
const onOpenChange = (value: boolean) => {
  emit("update:open", value)
}

// 创建项目
const createProject = useCreateProject({
  onSuccess: () => {
    toast.success(`项目 ${formData.name} 成功`)
    onOpenChange(false)
  },
  onError: () => {
    toast.error("创建失败")
  }
})

// 表单提交处理
const handleSubmit = () => {
  // 校验数据
  if (!formData.name || !formData.managerID || !formData.priority) {
    toast.error("请填写项目名称、负责人和优先级")
    return
  }

  const req: CreateProjectRequest = {
    name: formData.name,
    description: formData.description,
    managerID: formData.managerID,
    remark: formData.remark,
    priority: formData.priority,
    duration: formData.duration
  }

  createProject.mutate(req)
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
