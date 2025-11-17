import { listProcesses } from "@/api/projects"
import type { Process } from "@/model/processes"
import { useQuery } from "@tanstack/react-query"
import {
  getCoreRowModel,
  useReactTable,
  type ColumnDef
} from "@tanstack/react-table"
import { useParams } from "react-router"
import { toast } from "sonner"

const ProcessPage = () => {
  const params = useParams<{ id: string }>()
  const { data, isError, isPending } = useQuery({
    queryKey: ["processes", params.id],
    queryFn: () => {
      if (params.id) {
        return listProcesses(parseInt(params.id))
      }
      return []
    }
  })

  const columns: ColumnDef<Process>[] = [
    {
      accessorKey: "sequence",
      header: "序号"
    },
    {
      accessorKey: "processName",
      header: "流程"
    },
    {
      accessorKey: "description",
      header: "描述"
    },
    {
      accessorKey: "employeeName",
      header: "负责人"
    },
    {
      accessorKey: "status",
      header: "状态"
    },
    {
      header: "完成时间",
      cell: ({ row }) => {
        const startTime = row.original.startTime
        const endTime = row.original.endTime

        if (!endTime.Valid) {
          return "-"
        }

        const endDate = new Date(endTime.Time)
        const endDateFormat = `${endDate.getFullYear()}-${
          endDate.getMonth() + 1
        }-${endDate.getDate()} ${endDate.getHours()}:${endDate.getMinutes()}`
        if (!startTime.Valid) {
          return endDateFormat
        }

        const startDate = new Date(startTime.Time)
        const timeDiff = endDate.getTime() - startDate.getTime()
        const diffDays = Math.floor((timeDiff / 1000) * 3600 * 24)
        const diffHours = Math.floor(
          ((timeDiff % (1000 * 3600 * 24)) / 1000) * 3600
        )
        const diffMinutes = Math.floor(((timeDiff % (1000 * 3600)) / 1000) * 60)
        return (
          <div title={`耗时${diffDays}天${diffHours}小时${diffMinutes}分`}>
            {endDateFormat}
          </div>
        )
      }
    },
    {
      accessorKey: "remark",
      header: "备注"
    }
  ]

  const table = useReactTable({
    data: data || [],
    columns,
    getCoreRowModel: getCoreRowModel()
  })

  if (isError) {
    toast.error("获取流程失败")
    return null
  }

  if (isPending) {
    return <div>Loading...</div>
  }
  console.log(table)
  return <div>当前项目 ID: {params.id}</div>
}
export default ProcessPage
