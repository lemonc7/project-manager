<template>
  <Dialog :open="props.open" @update:open="onOpenChange">
    <DialogContent>
      <DialogTitle>编辑流程</DialogTitle>
      <DialogDescription>
        <span v-if="props.process && props.process.startedAt"
          >开始时间：{{
            new Date(props.process.startedAt).toISOString().split("T")[0]
          }}</span
        >
        <span v-if="props.process && props.process.finishedAt" class="ml-2">
          结束时间：{{
            new Date(props.process.finishedAt).toISOString().split("T")[0]
          }}
        </span>
      </DialogDescription>
      <div class="w-full max-w-md">
        <form @submit.prevent="handleSubmit">
          <FieldGroup>
            <FieldSet>
              <!-- 流程名称 -->
              <Field>
                <FieldLabel for="process-name">名称</FieldLabel>
                <Input
                  id="process-name"
                  v-model="formData.name"
                  placeholder="请输入流程名称"
                  required
                />
              </Field>
              <!-- 细节描述 -->
              <Field>
                <FieldLabel for="process-description">描述</FieldLabel>
                <Input
                  id="process-description"
                  v-model="formData.description"
                  placeholder="请输入细节描述"
                />
              </Field>
              <div class="grid grid-cols-3 gap-4">
                <!-- 负责人 -->
                <Field>
                  <FieldLabel>负责人</FieldLabel>
                  <EmployeesSelect v-model:employee-id="formData.employeeID" />
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
                        v-for="status in ProcessStatusOptions"
                        :key="status"
                        :value="status"
                      >
                        {{ status }}
                      </SelectItem>
                    </SelectContent>
                  </Select>
                </Field>
                <!-- 机时 -->
                <Field>
                  <FieldLabel for="process-spend-time">机时</FieldLabel>
                  <Input
                    id="process-spend-time"
                    v-model="formData.spendTime"
                    placeholder="请输入机时"
                    type="number"
                    min="0"
                    step="any"
                  />
                </Field>
              </div>
              <!-- 备注 -->
              <Field>
                <FieldLabel for="process-remark">备注</FieldLabel>
                <Textarea
                  id="process-remark"
                  placeholder="请输入流程备注"
                  class="resize-none"
                  v-model="formData.remark"
                />
              </Field>
            </FieldSet>
            <Field orientation="horizontal" class="justify-end">
              <Button
                type="button"
                variant="outline"
                @click="onOpenChange(false)"
              >
                取消
              </Button>
              <Button type="submit" :disabled="updateProcess.isPending.value">
                {{ updateProcess.isPending.value ? "更新中..." : "更新" }}
              </Button>
            </Field>
          </FieldGroup>
        </form>
      </div>
    </DialogContent>
  </Dialog>
</template>

<script setup lang="ts">
import useUpdateProcess from "@/composables/use-update-process"
import type { Employee } from "@/models/employees"
import {
  ProcessStatusOptions,
  type Process,
  type ProcessStatus,
  type UpdateProcessRequest
} from "@/models/processes"
import { reactive, watch } from "vue"
import { toast } from "vue-sonner"

const props = defineProps<{
  open: boolean
  employees: Employee[]
  process: Process | null
}>()

const emit = defineEmits<{
  (e: "update:open", value: boolean): void
}>()

// 默认值
const formData = reactive<{
  name: string
  description?: string
  remark?: string
  status: ProcessStatus
  employeeID?: number
  spendTime?: number
}>({
  name: "",
  status: "等待"
})

// 触发update:open事件，通知父组件关闭弹窗
const onOpenChange = (value: boolean) => {
  emit("update:open", value)
}

// 更新process
const updateProcess = useUpdateProcess({
  onSuccess: () => {
    toast.success(`流程 ${formData.name} 更新成功`)
    onOpenChange(false)
  },
  onError: () => {
    toast.error("更新失败")
  }
})

// 监听props.process变化
watch(
  () => props.process,
  (newProcess) => {
    if (newProcess) {
      Object.assign(formData, {
        name: newProcess.name,
        description: newProcess.description,
        remark: newProcess.remark,
        status: newProcess.status,
        employeeID: newProcess.employeeID,
        spendTime: newProcess.spendTime
      })
    }
  },
  { immediate: true }
)

// 监听open变化，用于取消操作后的数据重置
watch(
  () => props.open,
  (isOpen) => {
    if (!isOpen && props.process) {
      // 弹窗关闭时，将formData重置为原始的process值
      Object.assign(formData, {
        name: props.process.name,
        description: props.process.description,
        remark: props.process.remark,
        status: props.process.status,
        employeeID: props.process.employeeID,
        spendTime: props.process.spendTime
      })
    }
  }
)

// 表单提交处理
const handleSubmit = () => {
  if (!props.process?.id) return
  // 基础校验
  if (!formData.name || !formData.employeeID) {
    toast.error("请填写流程名称和负责人")
    return
  }

  // 构造请求
  const req: UpdateProcessRequest = {
    name: formData.name,
    description: formData.description,
    remark: formData.remark,
    status: formData.status,
    employeeID: formData.employeeID,
    spendTime: formData.spendTime
  }

  updateProcess.mutate({ id: props.process.id, req })
}
</script>
