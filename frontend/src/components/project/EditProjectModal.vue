<template>
  <UModal
    :open="props.open"
    @update:open="onOpenChange"
    title="编辑项目"
    description="更新成员的内容(不包括用户名和密码)"
  >
    <template #body>
      <UForm
        :schema="schema"
        :state="state"
        class="space-y-4"
        @submit="handleSubmit"
      >
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
import { reactive, watch } from "vue"
import type { Project, UpdateProjectRequest } from "../../models/projects"
import { useProjectStore } from "../../composables/useProjectStore"
import { storeToRefs } from "pinia"

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
const { updateProject } = projectStore
const { isUpdating } = storeToRefs(projectStore)

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

// 监听employee变化，用于数据初始化和外部数据更新
watch(
  () => props.project,
  (newVal) => {
    if (newVal) {
      Object.assign(state, {
        name: newVal.name,
        description: newVal.description,
        remark: newVal.remark,
        status: newVal.status,
        manager: newVal.managerID,
        priority: newVal.priority,
        duration: newVal.duration
      })
    }
  },
  { immediate: true }
)

// 监听open变化，关闭时重置表单
watch(
  () => props.open,
  (newVal) => {
    if (!newVal && props.project) {
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
