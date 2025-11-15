import type { Employee } from "@/model/employees"
import { create } from "zustand"

interface EmployeeList {
  employeeList: Employee[]
  setEmployeeList: (employeeList: Employee[]) => void
}

export const useEmployeeList = create<EmployeeList>((set) => ({
  employeeList: [],
  setEmployeeList: (employeeList) => set({ employeeList })
}))
