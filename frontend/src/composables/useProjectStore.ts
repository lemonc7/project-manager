import { useMutation, useQuery, useQueryClient } from "@tanstack/vue-query"
import { defineStore } from "pinia"
import type {
  CreateProjectRequest,
  ListProjectsRequest,
  UpdateProjectRequest
} from "../models/projects"
import { computed, ref } from "vue"
import {
  createProject,
  deleteProject,
  finishProject,
  getProject,
  listProjects,
  startProject,
  updateProject
} from "../api/projects"
import { useRoute } from "vue-router"

function buildProjectsQueryKey(req: ListProjectsRequest) {
  return ["projects", "overview", req.startAt ?? null, req.endAt ?? null]
}

export const useProjectStore = defineStore("project", () => {
  const queryClient = useQueryClient()

  // 查询数据
  const req = ref<ListProjectsRequest>({})

  const queryKey = computed(() => {
    return buildProjectsQueryKey(req.value)
  })

  function invalidateProjects() {
    queryClient.invalidateQueries({
      queryKey: ["projects"]
    })
  }

  const projectQuery = useQuery({
    queryKey: queryKey,
    queryFn: () => listProjects(req.value)
  })

  const route = useRoute()
  const projectID = computed(() => parseInt(route.params.id as string))
  const getProjectQuery = useQuery({
    queryKey: ["projects", projectID],
    queryFn: () => {
      const id = projectID.value
      if (!id || isNaN(id)) {
        throw new Error("缺少项目ID")
      }
      return getProject(id)
    },
    enabled: computed(() => !!projectID.value && !isNaN(projectID.value))
  })

  // 更改查询的key
  function setQuery(newReq: ListProjectsRequest) {
    req.value = newReq
  }

  // actions
  const createMutations = useMutation({
    mutationFn: (req: CreateProjectRequest) => createProject(req),
    onSuccess: () => invalidateProjects()
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, req }: { id: number; req: UpdateProjectRequest }) =>
      updateProject(id, req),
    onSuccess: () => invalidateProjects()
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => deleteProject(id),
    onSuccess: () => invalidateProjects()
  })

  const startMutation = useMutation({
    mutationFn: (id: number) => startProject(id),
    onSuccess: () => invalidateProjects()
  })

  const finishMutation = useMutation({
    mutationFn: (id: number) => finishProject(id),
    onSuccess: () => invalidateProjects()
  })

  return {
    projects: projectQuery.data,
    isLoading: projectQuery.isLoading,
    isError: projectQuery.isError,
    gotProject: getProjectQuery.data,
    isGetting: getProjectQuery.isLoading,
    isGettingError: getProjectQuery.isError,

    isCreating: createMutations.isPending,
    isUpdating: updateMutation.isPending,
    isDeleting: deleteMutation.isPending,
    isFinishing: finishMutation.isPending,
    isStarting: startMutation.isPending,

    setQuery,

    createProject: createMutations.mutateAsync,
    updateProject: updateMutation.mutateAsync,
    deleteProject: deleteMutation.mutateAsync,
    finishProject: finishMutation.mutateAsync,
    startProject: startMutation.mutateAsync
  }
})
