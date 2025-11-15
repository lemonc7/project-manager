import type { Project } from "@/model/projects"
import { create } from "zustand"

interface ProjectList {
  projectList: Project[]
  setProjectList: (projectList: Project[]) => void
}

export const useProjectList = create<ProjectList>((set) => ({
  projectList: [],
  setProjectList: (projectList) => set({ projectList })
}))
