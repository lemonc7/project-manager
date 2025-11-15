import { getProcessesByProjectID } from "@/api/processes"
import type { Process } from "@/model/processes"
import { create } from "zustand"

interface ProcessCacheItem {
  processes: Process[]
}

interface ProcessStore {
  processCache: Record<number, ProcessCacheItem>
  fetchProcesses: (projectID: number, force?: boolean) => Promise<void>
  clearCache: () => void
}

export const useProcessStore = create<ProcessStore>((set, get) => ({
  processCache: {},
  // 缓存优先
  fetchProcesses: async (projectID, force = false) => {
    const cache = get().processCache[projectID]
    // 缓存存在且不是强制刷新，直接返回
    if (cache && !force) {
      return
    }
    // 从后端获取数据
    const data = await getProcessesByProjectID(projectID)
    // 放入缓存
    set((state) => ({
      processCache: {
        ...state.processCache,
        [projectID]: {
          processes: data
        }
      }
    }))
  },
  clearCache: () => {
    set({
      processCache: {}
    })
  }
}))
