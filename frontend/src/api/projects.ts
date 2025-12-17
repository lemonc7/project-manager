import type {
  CopyProjectRequest,
  CreateProjectRequest,
  ListProjectsRequest,
  Project,
  ProjectGet,
  UpdateProjectRequest
} from "../models/projects"
import service from "./request"

export async function listProjects(req: ListProjectsRequest) {
  const res = await service.post<Project[]>(`/project/list`, req)
  return Array.isArray(res.data) ? res.data : []
}

export async function createProject(req: CreateProjectRequest) {
  return service.post("/project", req)
}

export async function updateProject(id: number, req: UpdateProjectRequest) {
  return service.patch(`/project/${id}`, req)
}

export async function deleteProject(id: number) {
  return service.delete(`/project/${id}`)
}

export async function startProject(id: number) {
  return service.post(`/project/start/${id}`)
}

export async function finishProject(id: number) {
  return updateProject(id, { status: "完成" })
}

export async function getProject(id: number) {
  const res = await service.get<ProjectGet>(`/project/${id}`)
  return res.data
}

export async function resetProject(id: number) {
  return service.post(`/project/reset/${id}`)
}

export async function copyProject(id: number, req: CopyProjectRequest) {
  return service.post(`/project/copy/${id}`, req)
}

export async function isProjectExists(name: string) {
  const res = await service.get<boolean>("/project/exist", {
    params: {
      newName: name
    }
  })
  return res.data
}

export async function exportProject(id: number, name: string, toast: any) {
  try {
    const res = await service.get(`/project/export/${id}`, {
      responseType: "blob"
    })

    // 检查响应是不是文件流
    if (res.data.type === "application/json") {
      const text = await res.data.text()
      const errData = JSON.parse(text)
      toast.add({
        title: "导出失败",
        description: errData.message || "服务器内部错误",
        color: "error"
      })
    }

    // 创建文件下载
    const blob = new Blob([res.data], { type: "text/csv;charset=utf-8;" })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement("a")
    link.href = url

    // 使用项目名称作为文件名
    const dateStr = new Date().toISOString().split("T")[0]
    const filename = `${name}_${dateStr}.csv`
    link.setAttribute("download", filename)

    document.body.appendChild(link)
    link.click()

    // 善后处理
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)

    toast.add({
      title: "导出成功",
      description: `文件 ${filename} 已下载`,
      color: "success"
    })
  } catch (err) {
    console.error("导出失败:", err)
    toast.add({
      title: "导出失败",
      description: "服务器内部错误",
      color: "error"
    })
  }
}
