import { listProjectsFilter } from "@/api/projects"
import type { ListProjectsRequest } from "@/model/projects"
import { useQuery } from "@tanstack/react-query"

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

export default function useProjects(req: ListProjectsRequest) {
  return useQuery({
    queryKey: buildProjectsQueryKey(req),
    queryFn: () => listProjectsFilter(req)
  })
}
