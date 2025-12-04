import { listProjects } from "../api/projects"
import type { FilterStatus, ListProjectsRequest } from "../models/projects"
import { useQuery } from "@tanstack/vue-query"

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

export default function useProjects({
  filter,
  req
}: {
  filter: FilterStatus
  req: ListProjectsRequest
}) {
  return useQuery({
    queryKey: buildProjectsQueryKey(filter, req),
    queryFn: () => listProjects(filter, req)
  })
}
