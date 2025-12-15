import { defineStore } from "pinia"
import { computed, ref } from "vue"
import type { Employee } from "../models/employees"

export const useEmployeeStore = defineStore("employee", () => {
  const employeeMap = ref(new Map<number, string>())

  function setEmployees(list: Employee[]) {
    const map = new Map<number, string>()
    list.forEach((e) => map.set(e.id, e.name))
    employeeMap.value = map
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
    employeeItems,
    getEmployeeName,
    setEmployees
  }
})
