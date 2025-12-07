import { computed, ref, watch } from "vue"
import useEmployees from "./useEmployees"

// 将状态放在外部，保持单例
const employeeMap = ref<Map<number, string>>(new Map())
let isInitialized = false

export default function useEmployeeMap() {
  // 只有在第一次调用才初始化监听
  if (!isInitialized) {
    const { data } = useEmployees()
    watch(
      data,
      (list) => {
        const newMap = new Map<number, string>()
        list?.forEach((e) => newMap.set(e.id, e.name))
        employeeMap.value = newMap
      },
      { immediate: true }
    )
    isInitialized = true
  }

  function getEmployeeName(id?: number) {
    if (id === undefined) return "未分配"
    return employeeMap.value.get(id) ?? "未知"
  }

  const employeeItems = computed(() =>
    Array.from(employeeMap.value.entries(), ([id, name]) => ({
      value: id,
      label: name
    }))
  )

  return {
    employeeMap,
    getEmployeeName,
    employeeItems
  }
}
