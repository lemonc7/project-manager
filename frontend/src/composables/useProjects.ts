import { useMutation, useQuery, useQueryClient } from "@tanstack/vue-query"
import type {
  CopyProjectRequest,
  CreateProjectRequest,
  ListProjectsRequest,
  UpdateProjectRequest
} from "../models/projects"
import { computed, ref } from "vue"
import {
  copyProject,
  createProject,
  deleteProject,
  finishProject,
  getProject,
  isProjectExists,
  listProjects,
  resetProject,
  startProject,
  updateProject
} from "../api/projects"
import { useRoute } from "vue-router"

function buildProjectsQueryKey(req: ListProjectsRequest) {
  return [
    "projects",
    "overview",
    req.managerID ?? null,
    req.status ?? null,
    req.startAt ?? null,
    req.endAt ?? null
  ]
}

export function useProjects() {
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
    queryFn: () => listProjects(req.value),
    refetchInterval: 1000 * 60,
    staleTime: Infinity
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

  const resetMutation = useMutation({
    mutationFn: (id: number) => resetProject(id),
    onSuccess: () => invalidateProjects()
  })

  const copyMutation = useMutation({
    mutationFn: ({ id, req }: { id: number; req: CopyProjectRequest }) =>
      copyProject(id, req),
    onSuccess: () => invalidateProjects()
  })

  const existMutation = useMutation({
    mutationFn: (name: string) => isProjectExists(name)
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
    isResetting: resetMutation.isPending,
    isCopying: copyMutation.isPending,
    isChecking: existMutation.isPending,

    setQuery: (request: ListProjectsRequest) => {
      if (request.managerID === 0) {
        request.managerID = undefined
      }
      req.value = request
    },

    createProject: createMutations.mutateAsync,
    updateProject: updateMutation.mutateAsync,
    deleteProject: deleteMutation.mutateAsync,
    finishProject: finishMutation.mutateAsync,
    startProject: startMutation.mutateAsync,
    resetProject: resetMutation.mutateAsync,
    copyProject: copyMutation.mutateAsync,
    checkProject: existMutation.mutateAsync
  }
}
