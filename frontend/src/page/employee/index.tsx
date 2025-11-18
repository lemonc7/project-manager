import { listEmployees } from "@/api/employees"
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from "@/components/ui/table"
import type { Employee } from "@/model/employees"
import { useQuery } from "@tanstack/react-query"
import { toast } from "sonner"

function EmployeePage() {
  const { data, isError, isPending } = useQuery<Employee[], Error>({
    queryKey: ["employees"],
    queryFn: listEmployees
  })

  if (isError) {
    toast.error("获取员工失败")
    return null
  }

  if (isPending) {
    return <div>Loading...</div>
  }

  return (
    <Table>
      <TableCaption>员工表</TableCaption>
      <TableHeader>
        <TableRow>
          <TableHead>员工</TableHead>
          <TableHead>小组</TableHead>
          <TableHead>添加时间</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {data.map(
          (emp) =>
            emp.isActive && (
              <TableRow key={emp.id}>
                <TableCell>{emp.employeeName}</TableCell>
                <TableCell>{emp.teamName}</TableCell>
                <TableCell>
                  {emp.createdAt && formatDate(emp.createdAt)}
                </TableCell>
              </TableRow>
            )
        )}
      </TableBody>
    </Table>
  )
}

function formatDate(dateStr: string) {
  const date = new Date(dateStr)
  const year = date.getFullYear()
  const month = date.getMonth() + 1
  const day = date.getDate()
  return `${year}-${month}-${day}`
}

export default EmployeePage
