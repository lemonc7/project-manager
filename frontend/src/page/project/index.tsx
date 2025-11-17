import { listProjects } from "@/api/projects"
import type { Project } from "@/model/projects"
import { useQuery } from "@tanstack/react-query"
import { toast } from "sonner"
import { type ColumnDef } from "@tanstack/react-table"

const ProjectPage = () => {
  const { data, isError, isPending } = useQuery<Project[], Error>({
    queryKey: ["projects", "overview"],
    queryFn: listProjects
  })

  if (isPending) {
    return <div>Loading...</div>
  }

  if (isError) {
    toast.error("获取项目失败")
    return null
  }

  const columns: ColumnDef<Project>[] = []

  return (
    <div>
      <h1>项目总览</h1>
    </div>
  )
}
export default ProjectPage
