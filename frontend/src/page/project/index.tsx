import { listProjects } from "@/api/projects"
import type { Project } from "@/model/projects"
import { useQuery } from "@tanstack/react-query"
import { toast } from "sonner"
import {
  flexRender,
  getCoreRowModel,
  useReactTable,
  type ColumnDef
} from "@tanstack/react-table"
import { useNavigate } from "react-router"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from "@/components/ui/table"

const ProjectPage = () => {
  const { data, isError, isPending } = useQuery<Project[], Error>({
    queryKey: ["projects", "overview"],
    queryFn: listProjects
  })
  const navigate = useNavigate()

  const columns: ColumnDef<Project>[] = [
    {
      accessorKey: "projectName",
      header: "项目",
      cell: ({ row }) => (
        <div
          className="cursor-pointer"
          onClick={() => navigate(`/projects/${row.original.id}`)}
        >
          {row.getValue("projectName")}
        </div>
      )
    },
    {
      accessorKey: "description",
      header: "描述"
    },
    {
      accessorKey: "managerName",
      header: "负责人"
    },
    {
      header: "剩余时间",
      cell: ({ row }) => {
        const startTime = row.original.startTime
        const duration = row.original.duration

        // 没设置startTime或duration为0时，说明没设置时间
        if (!startTime.Valid || duration === 0) {
          return "-"
        }
        // 计算剩余天数
        const startDate = new Date(startTime.Time)
        const currentDate = new Date()

        const endDate = new Date(startTime.Time)
        endDate.setDate(startDate.getDate() + duration)

        const timeDiff = endDate.getTime() - currentDate.getTime()
        const remainingDays = Math.ceil(timeDiff / (1000 * 3600 * 24))

        return remainingDays >= 0
          ? `剩余${remainingDays}天`
          : `超时${-remainingDays}天`
      }
    },
    {
      accessorKey: "status",
      header: "状态"
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

  if (isPending) {
    return <div>Loading...</div>
  }

  if (isError) {
    toast.error("获取项目失败")
    return null
  }
  return (
    <div className="overflow-hidden rounded-md border">
      <Table>
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => {
                return (
                  <TableHead key={header.id}>
                    {header.isPlaceholder
                      ? null
                      : flexRender(
                          header.column.columnDef.header,
                          header.getContext()
                        )}
                  </TableHead>
                )
              })}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {table.getRowModel().rows?.length ? (
            table.getRowModel().rows.map((row) => (
              <TableRow
                key={row.id}
                data-state={row.getIsSelected() && "selected"}
              >
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={columns.length} className="h-24 text-center">
                无结果
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  )
}
export default ProjectPage
