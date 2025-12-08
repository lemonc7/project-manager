import { useMutation, useQuery, useQueryClient } from "@tanstack/vue-query"
import { defineStore } from "pinia"
import type {
  CreateProjectRequest,
  FilterStatus,
  ListProjectsRequest,
  UpdateProjectRequest
} from "../models/projects"
import { computed, ref } from "vue"
import {
  createProject,
  deleteProject,
  finishProject,
  listProjects,
  startProject,
  updateProject
} from "../api/projects"

function buildProjectsQueryKey(filter: FilterStatus, req: ListProjectsRequest) {
  return [
    "projects",
    "overview",
    filter,
    req.managerID ?? null,
    req.status ?? null,
    req.startAt ?? null,
    req.endAt ?? null
  ]
}

export const useProjectStore = defineStore("project", () => {
  const queryClient = useQueryClient()

  // 查询数据
  const filter = ref<FilterStatus>("filter")
  const req = ref<ListProjectsRequest>({})

  const queryKey = computed(() => {
    return buildProjectsQueryKey(filter.value, req.value)
  })

  const projectQuery = useQuery({
    queryKey: queryKey,
    queryFn: () => listProjects(filter.value, req.value)
  })

  // 更改查询的key
  function setQuery(newFilter: FilterStatus, newReq: ListProjectsRequest) {
    filter.value = newFilter
    req.value = newReq
  }

  // actions
  const createMutations = useMutation({
    mutationFn: (req: CreateProjectRequest) => createProject(req),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects", "overview"] })
    }
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, req }: { id: number; req: UpdateProjectRequest }) =>
      updateProject(id, req),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects", "overview"] })
    }
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => deleteProject(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects", "overview"] })
    }
  })

  const startMutation = useMutation({
    mutationFn: (id: number) => startProject(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects", "overview"] })
    }
  })

  const finishMutation = useMutation({
    mutationFn: (id: number) => finishProject(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects", "overview"] })
    }
  })

  return {
    projects: projectQuery.data,
    isLoading: projectQuery.isLoading,
    isError: projectQuery.isError,

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
