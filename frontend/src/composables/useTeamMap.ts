import { computed, ref, watch } from "vue"
import useTeams from "./useTeams"

// 将状态放在外部，保持单例
const teamMap = ref<Map<number, string>>(new Map())
let isInitialized = false

export default function useTeamMap() {
  // 只有在第一次调用 useTeamMap 时才初始化监听
  if (!isInitialized) {
    const { data } = useTeams()
    watch(
      data,
      (list) => {
        const newMap = new Map<number, string>()
        list?.forEach((t) => newMap.set(t.id, t.name))
        teamMap.value = newMap
      },
      { immediate: true }
    )
    isInitialized = true
  }

  function getTeamName(id?: number) {
    if (id === undefined) return "未分配"
    return teamMap.value.get(id) ?? "未知"
  }

  const teamItems = computed(() =>
    Array.from(teamMap.value.entries(), ([id, name]) => ({
      value: id,
      label: name
    }))
  )

  return {
    teamMap,
    getTeamName,
    teamItems
  }
}
