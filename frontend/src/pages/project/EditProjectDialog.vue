<template>
  <Dialog :open="props.open" @update:open="onOpenChange">
    <DialogContent>
      <DialogTitle>编辑项目</DialogTitle>
      <DialogDescription>
        <span
          v-if="
            !props.project ||
            props.project.remainingDays === undefined ||
            !props.project.startedAt
          "
          >项目未开始或设置周期</span
        >
        <span v-else>
          开始时间：{{
            new Date(props.project.startedAt).toISOString().split("T")[0]
          }}{{
            props.project.remainingDays >= 0
              ? `(剩余${props.project.remainingDays}天)`
              : `(超时${-props.project.remainingDays}天)`
          }}
        </span>
      </DialogDescription>
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
              <div class="grid grid-cols-2 gap-4">
                <!-- 负责人 -->
                <Field>
                  <FieldLabel>负责人</FieldLabel>
                  <EmployeesSelect v-model:employee-id="formData.managerID" />
                </Field>
                <!-- 状态 -->
                <Field>
                  <FieldLabel>状态</FieldLabel>
                  <Select v-model="formData.status">
                    <SelectTrigger>
                      <SelectValue placeholder="选择状态" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem
                        v-for="status in ProjectStatusOptions"
                        :key="status"
                        :value="status"
                      >
                        {{ status }}
                      </SelectItem>
                    </SelectContent>
                  </Select>
                </Field>
              </div>
              <div class="grid grid-cols-2 gap-4">
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
              <Button type="submit" :disabled="updateProject.isPending.value">
                {{ updateProject.isPending.value ? "更新中..." : "更新" }}
              </Button>
            </Field>
          </FieldGroup>
        </form>
      </div>
    </DialogContent>
  </Dialog>
</template>

<script setup lang="ts">
import useUpdateProject from "@/composables/use-update-project"
import type { Employee } from "@/models/employees"
import {
  ProjectStatusOptions,
  type Project,
  type ProjectStatus,
  type UpdateProjectRequest
} from "@/models/projects"
import { reactive, watch } from "vue"
import { toast } from "vue-sonner"

const props = defineProps<{
  open: boolean
  employees: Employee[]
  project: Project | null
}>()

const emit = defineEmits<{
  (e: "update:open", value: boolean): void
}>()

const formData = reactive<{
  name: string
  description?: string
  remark?: string
  status: ProjectStatus
  managerID: number
  priority: number
  duration?: number
}>({
  name: "",
  status: "排单",
  managerID: 1,
  priority: 5
})

const onOpenChange = (value: boolean) => {
  emit("update:open", value)
}

const updateProject = useUpdateProject({
  onSuccess: () => {
    toast.success(`项目 ${formData.name} 成功`)
    onOpenChange(false)
  },
  onError: () => {
    toast.error("更新失败")
  }
})

// 监听props.project变化，用于数据初始化和外部数据更新
watch(
  () => props.project,
  (newProject) => {
    if (newProject) {
      Object.assign(formData, {
        name: newProject.name,
        description: newProject.description,
        remark: newProject.remark,
        status: newProject.status,
        managerID: newProject.managerID,
        priority: newProject.priority,
        startedAt: newProject.startedAt,
        duration: newProject.duration,
        finishedAt: newProject.finishedAt
      })
    }
  },
  { immediate: true }
)
// 监听open变化，用于取消操作后的数据重置
watch(
  () => props.open,
  (isOpen) => {
    if (!isOpen && props.project) {
      // 弹窗关闭时，将formData重置为原始的project值
      Object.assign(formData, {
        name: props.project.name,
        description: props.project.description,
        remark: props.project.remark,
        status: props.project.status,
        managerID: props.project.managerID,
        priority: props.project.priority,
        startedAt: props.project.startedAt,
        duration: props.project.duration,
        finishedAt: props.project.finishedAt
      })
    }
  }
)

// 表单提交处理
const handleSubmit = () => {
  if (!props.project?.id) return
  // 基础校验
  if (!formData.name || !formData.managerID || !formData.priority) {
    toast.error("请填写项目名称、负责人和优先级")
    return
  }

  // 构造请求
  const req: UpdateProjectRequest = {
    name: formData.name,
    description: formData.description,
    remark: formData.remark,
    status: formData.status,
    managerID: formData.managerID,
    priority: formData.priority,
    duration: formData.duration
  }

  updateProject.mutate({ id: props.project.id, req })
}
</script>
