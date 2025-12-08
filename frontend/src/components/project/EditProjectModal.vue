<template>
  <UModal
    :open="props.open"
    @update:open="onOpenChange"
    title="编辑"
    :description="detail"
  >
    <template #body>
      <UForm
        :schema="schema"
        :state="state"
        class="space-y-4"
        @submit="handleSubmit"
      >
        <!-- 项目名称 -->
        <UFormField label="名称" name="name" required>
          <UInput
            v-model="state.name"
            class="w-full"
            placeholder="请输入项目名称"
          />
        </UFormField>
        <!-- 细节描述 -->
        <UFormField label="细节描述" name="description">
          <UInput
            v-model="state.description"
            class="w-full"
            placeholder="请输入细节描述"
          />
        </UFormField>
        <!-- 负责人+状态 -->
        <div class="grid grid-cols-2 gap-4">
          <UFormField label="负责人" name="manager">
            <USelectMenu
              v-model="state.manager"
              :items="employeeItems"
              value-key="value"
              placeholder="请选择负责人"
              class="w-full"
            />
          </UFormField>
          <UFormField label="状态" name="status">
            <USelect
              v-model="state.status"
              :items="[...ProjectStatusOptions]"
              placeholder="请选择状态"
              class="w-full"
            />
          </UFormField>
        </div>
        <!-- 优先级+项目周期 -->
        <div class="grid grid-cols-2 gap-2">
          <UFormField label="优先级" name="priority" required>
            <UInput
              type="number"
              v-model="state.priority"
              class="w-full"
              placeholder="请输入优先级"
            />
          </UFormField>
          <UFormField label="项目周期" name="duration">
            <UInput
              type="number"
              v-model="state.duration"
              class="w-full"
              placeholder="请输入项目周期"
            />
          </UFormField>
        </div>
        <!-- 备注 -->
        <UFormField>
          <UTextarea
            v-model="state.remark"
            class="w-full"
            placeholder="请输入备注"
            autoresize
          />
        </UFormField>
        <!-- 按钮 -->
        <div class="flex justify-end gap-2">
          <UButton
            label="取消"
            variant="outline"
            @click="onOpenChange(false)"
          />
          <UButton label="更新" type="submit" :loading="isUpdating" />
        </div>
      </UForm>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import z from "zod"
import { computed, reactive, watch } from "vue"
import {
  ProjectStatusOptions,
  type Project,
  type UpdateProjectRequest
} from "../../models/projects"
import { useProjectStore } from "../../composables/useProjectStore"
import { storeToRefs } from "pinia"
import { useEmployeeStore } from "../../composables/useEmployeeStore"

const props = defineProps<{
  open: boolean
  project: Project | null
}>()

const emit = defineEmits<{
  (e: "update:open", open: boolean): void
}>()

const toast = useToast()

// 定义表单数据
const schema = z.object({
  name: z.string().nonempty("请输入姓名"),
  description: z.string().optional(),
  remark: z.string().optional(),
  status: z.enum(["排单", "进行", "暂停", "完成", "终止"]),
  manager: z.number(),
  priority: z.number().min(1, "优先级不能小于1").max(9, "优先级不能大于9"),
  duration: z.number().min(0, "项目周期不能小于0").optional()
})

type Schema = z.output<typeof schema>

const state = reactive<Schema>({
  name: "",
  description: undefined,
  remark: undefined,
  status: "排单",
  manager: 0,
  priority: 5,
  duration: undefined
})

function onOpenChange(open: boolean) {
  emit("update:open", open)
}

const projectStore = useProjectStore()
const employeeStore = useEmployeeStore()
const { updateProject } = projectStore
const { isUpdating } = storeToRefs(projectStore)
const { employeeItems } = storeToRefs(employeeStore)

async function handleSubmit() {
  if (!props.project?.id) return
  const req: UpdateProjectRequest = {
    name: state.name,
    description: state.description,
    remark: state.remark,
    status: state.status,
    managerID: state.manager,
    priority: state.priority,
    duration: state.duration
  }
  try {
    await updateProject({ id: props.project.id, req })
    toast.add({
      title: `项目 ${state.name} 更新成功`,
      color: "success"
    })
    onOpenChange(false)
  } catch (err) {
    toast.add({
      title: "更新失败",
      color: "error"
    })
    console.log("更新项目失败:", err)
  }
}

const detail = computed(() => {
  const p = props.project
  if (!p?.startedAt) return "项目未开始"
  const date = p.startedAt.slice(0, 10)
  let remainingDays = ""
  if (p.remainingDays === undefined) {
    remainingDays = "(项目周期未设置)"
  } else if (p.remainingDays >= 0) {
    remainingDays = `(剩余${p.remainingDays}天)`
  } else {
    remainingDays = `(超时${Math.abs(p.remainingDays)}天)`
  }

  return `项目开始时间:${date} ${remainingDays}`
})

// 监听open变化，打开时同步表单数据
watch(
  () => props.open,
  (isOpen) => {
    if (isOpen && props.project) {
      Object.assign(state, {
        name: props.project.name,
        description: props.project.description,
        remark: props.project.remark,
        status: props.project.status,
        manager: props.project.managerID,
        priority: props.project.priority,
        duration: props.project.duration
      })
    }
  }
)
</script>
