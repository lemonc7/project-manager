<template>
  <!-- Header中的按钮 -->
  <SafeTeleport to="#header-actions">
    <Input v-model="filterText" placeholder="项目名称过滤..." />
    <Button @click="isCreateDialogOpen = true">创建项目</Button>
  </SafeTeleport>
  <!-- 表格渲染 -->
  <div>
    <Tabs default-value="doing">
      <TabsList>
        <TabsTrigger
          v-for="item in tabItems"
          :key="item.value"
          :value="item.value"
          class="w-16"
        >
          {{ item.title }}
        </TabsTrigger>
      </TabsList>
      <TabsContent
        v-for="item in tabItems"
        :key="item.value"
        :value="item.value"
      >
        <ProjectTable
          class="py-2"
          :project-filter="item.value"
          :request="{}"
          :disable-select="!hasRole('manager', 'admin')"
          :filter-text="filterText"
          @edit="handleEditProject"
          @delete="handleDeleteProject"
          @start="handleStartProject"
          @finish="handleFinishProject"
        />
      </TabsContent>
    </Tabs>
  </div>
  <!-- 创建项目弹窗 -->
  <CreateProjectDialog
    v-model:open="isCreateDialogOpen"
    :employees="employees ?? []"
  />
  <!-- 编辑项目弹窗 -->
  <EditProjectDialog
    v-model:open="isEditDialogOpen"
    :employees="employees ?? []"
    :project="selectedProject"
  />
  <!-- 删除弹窗 -->
  <AlertDialog v-model:open="isDeleteDialogOpen">
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle> 删除项目 </AlertDialogTitle>
        <AlertDialogDescription>
          确定要删除
          <span class="font-bold text-red-500">{{
            selectedProject?.name
          }}</span>
          项目吗? 这个操作不能撤销，将从数据库中删除该项目以及相关工艺流程。
        </AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogCancel
          @click="selectedProject = null"
          :disabled="deleteProject.isPending.value"
        >
          取消
        </AlertDialogCancel>
        <AlertDialogAction
          @click="handleDeleteConfirm"
          :disabled="deleteProject.isPending.value"
        >
          {{ deleteProject.isPending.value ? "删除中..." : "确认" }}
        </AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
  <!-- 开始弹窗 -->
  <AlertDialog v-model:open="isStartDialogOpen">
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>开始项目</AlertDialogTitle>
        <AlertDialogDescription>
          确定要开始<span class="font-bold text-red-500">{{
            selectedProject?.name
          }}</span
          >项目吗?
          <span v-if="!selectedProject?.duration">
            <span class="font-bold text-red-500">项目周期</span>未设置，请确认！
          </span>
        </AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogCancel
          @click="selectedProject = null"
          :disabled="startProject.isPending.value"
        >
          取消
        </AlertDialogCancel>
        <AlertDialogAction
          @click="handleStartConfirm"
          :disabled="startProject.isPending.value"
        >
          {{ startProject.isPending.value ? "操作中..." : "确认" }}
        </AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
  <!-- 完成弹窗 -->
  <AlertDialog v-model:open="isFinishDialogOpen">
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>完成项目</AlertDialogTitle>
        <AlertDialogDescription>
          确定项目
          <span class="font-bold text-red-500">{{
            selectedProject?.name
          }}</span>
          已完成吗?
        </AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogCancel
          @click="selectedProject = null"
          :disabled="finishProject.isPending.value"
        >
          取消
        </AlertDialogCancel>
        <AlertDialogAction
          @click="handleFinishConfirm"
          :disabled="finishProject.isPending.value"
        >
          {{ finishProject.isPending.value ? "操作中..." : "确认" }}
        </AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>

<script setup lang="ts">
import useDeleteProject from "@/composables/use-delete-project"
import useEmployees from "@/composables/use-employees"
import useFinishProject from "@/composables/use-finish-project"
import useStartProject from "@/composables/use-start-project"
import type { Project } from "@/models/projects"
import { hasRole } from "@/util/auth"
import { ref } from "vue"
import { toast } from "vue-sonner"

const tabItems = [
  { title: "进行", value: "doing" },
  { title: "待办", value: "todo" },
  { title: "完成", value: "done" }
] as const

const { data: employees } = useEmployees()

const isCreateDialogOpen = ref(false)
const isEditDialogOpen = ref(false)
const isDeleteDialogOpen = ref(false)
const isStartDialogOpen = ref(false)
const isFinishDialogOpen = ref(false)

const selectedProject = ref<Project | null>(null)

const filterText = ref("")

const handleEditProject = (project: Project) => {
  selectedProject.value = project
  isEditDialogOpen.value = true
}

const handleDeleteProject = (project: Project) => {
  selectedProject.value = project
  isDeleteDialogOpen.value = true
}

const handleStartProject = (project: Project) => {
  selectedProject.value = project
  isStartDialogOpen.value = true
}

const handleFinishProject = (project: Project) => {
  selectedProject.value = project
  isFinishDialogOpen.value = true
}

const deleteProject = useDeleteProject({
  onSuccess: () => {
    toast.success(`项目 ${selectedProject.value?.name} 删除成功`)
    isDeleteDialogOpen.value = false
    selectedProject.value = null
  },
  onError: () => {
    isDeleteDialogOpen.value = false
    selectedProject.value = null
    toast.error("删除失败")
  }
})

const startProject = useStartProject({
  onSuccess: () => {
    toast.success(`项目 ${selectedProject.value?.name} 已开始`)
    isStartDialogOpen.value = false
    selectedProject.value = null
  },
  onError: () => {
    isStartDialogOpen.value = false
    selectedProject.value = null
    toast.error("操作失败")
  }
})

const finishProject = useFinishProject({
  onSuccess: () => {
    toast.success(`项目 ${selectedProject.value?.name} 已完成`)
    isFinishDialogOpen.value = false
    selectedProject.value = null
  },
  onError: () => {
    isFinishDialogOpen.value = false
    selectedProject.value = null
    toast.error("操作失败")
  }
})

const handleDeleteConfirm = (e: MouseEvent) => {
  e.preventDefault()
  if (selectedProject.value) {
    deleteProject.mutate(selectedProject.value.id)
  }
}

const handleStartConfirm = (e: MouseEvent) => {
  e.preventDefault()
  if (selectedProject.value) {
    startProject.mutate(selectedProject.value.id)
  }
}

const handleFinishConfirm = (e: MouseEvent) => {
  e.preventDefault()
  if (selectedProject.value) {
    finishProject.mutate(selectedProject.value.id)
  }
}
</script>
