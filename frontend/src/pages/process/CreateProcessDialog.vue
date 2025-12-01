<template>
  <Dialog :open="props.open" @update:open="onOpenChange">
    <DialogContent>
      <DialogTitle>创建流程</DialogTitle>
      <DialogDescription></DialogDescription>
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
              <div class="grid grid-cols-2 gap-4">
                <Field>
                  <FieldLabel>负责人</FieldLabel>
                  <EmployeesSelect v-model:employee-id="formData.employeeID" />
                </Field>
                <Field>
                  <FieldLabel for="process-sequence">流程序号</FieldLabel>
                  <Input
                    id="process-sequence"
                    v-model="formData.sequence"
                    placeholder="请输入流程序号"
                    type="number"
                    min="1"
                    :max="props.maxSequence"
                  />
                </Field>
              </div>
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
                variant="outline"
                type="button"
                @click="onOpenChange(false)"
              >
                取消
              </Button>
              <Button type="submit" :disabled="createProcess.isPending.value">
                {{ createProcess.isPending.value ? "创建中..." : "创建" }}
              </Button>
            </Field>
          </FieldGroup>
        </form>
      </div>
    </DialogContent>
  </Dialog>
</template>

<script setup lang="ts">
import useCreateProcess from "@/composables/use-create-process"
import type { Employee } from "@/models/employees"
import type { InsertProcessRequest } from "@/models/processes"
import { reactive, watch } from "vue"
import { toast } from "vue-sonner"

const props = defineProps<{
  open: boolean
  employees: Employee[]
  maxSequence: number
}>()

const emit = defineEmits<{
  (e: "update:open", value: boolean): void
}>()

const initialFormData = {
  name: "",
  sequence: undefined,
  employeeID: undefined,
  description: undefined,
  remark: undefined
}

const formData = reactive({ ...initialFormData })

// 触发update:open事件，通知父组件关闭弹窗
const onOpenChange = (value: boolean) => {
  emit("update:open", value)
}

// 创建流程
const createProcess = useCreateProcess({
  onSuccess: () => {
    toast.success(`流程 ${formData.name} 创建成功`)
    onOpenChange(false)
  },
  onError: () => {
    toast.error("创建失败")
  }
})

// 表单提交处理
const handleSubmit = () => {
  // 校验数据
  if (!formData.name) {
    toast.error("请填写流程名称")
    return
  }

  const req: Omit<InsertProcessRequest, "projectID"> = {
    name: formData.name,
    employeeID: formData.employeeID,
    description: formData.description,
    remark: formData.remark
  }

  createProcess.mutate({ sequence: formData.sequence, req })
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
