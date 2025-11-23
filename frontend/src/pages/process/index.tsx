import { updateProcess } from "@/api/processes"
import { listProcesses } from "@/api/projects"
import EditableCell from "@/components/common/editable-cell"
import { Loading } from "@/components/common/loading"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from "@/components/ui/table"
import type {
  Process,
  ProcessStatus,
  UpdateProcessRequest
} from "@/models/processes"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import {
  flexRender,
  getCoreRowModel,
  useReactTable,
  type ColumnDef
} from "@tanstack/react-table"
import { useParams } from "react-router"
import { toast } from "sonner"

const processStatusOptions = [
  "等待",
  "准备",
  "进行",
  "完成",
  "返工",
  "终止",
  "取消",
  "暂停"
] satisfies ProcessStatus[]

function ProcessPage() {
  const params = useParams<{ id: string }>()
  const { data, isError, isPending } = useQuery({
    queryKey: ["processes", params.id],
    queryFn: () => {
      if (params.id) {
        const api = listProcesses(+params.id)
        toast.promise(api, {
          loading: "流程加载中...",
          success: "流程加载成功",
          error: "加载失败"
        })
        return api
      }
      return []
    }
  })

  const queryClient = useQueryClient()

  const updateProcessMutation = useMutation({
    mutationFn: ({
      id,
      fields
    }: {
      id: number
      fields: Partial<UpdateProcessRequest>
    }) => updateProcess(id, fields),
    onSuccess: (updatedProcess) => {
      queryClient.setQueryData<Process[]>(
        ["processes", params.id],
        (currentProcesses) =>
          currentProcesses
            ? currentProcesses.map((process) =>
                process.id === updatedProcess.id
                  ? { ...process, ...updatedProcess }
                  : process
              )
            : [updatedProcess]
      )
    },
    onError: () => {
      toast.error("更新失败")
    }
  })

  const columns: ColumnDef<Process>[] = [
    {
      accessorKey: "sequence",
      header: "序号"
    },
    {
      header: "流程",
      cell: ({ row }) => {
        return (
          <EditableCell
            value={row.original.name}
            type="text"
            onSave={(value) =>
              updateProcessMutation.mutate({
                id: row.original.id,
                fields: { name: value }
              })
            }
          />
        )
      }
    },
    {
      header: "描述",
      cell: ({ row }) => {
        return (
          <EditableCell
            value={row.original.description}
            onSave={(value) =>
              updateProcessMutation.mutate({
                id: row.original.id,
                fields: { description: value }
              })
            }
          />
        )
      }
    },
    {
      accessorKey: "employeeName",
      header: "负责人"
    },
    {
      header: "状态",
      cell: ({ row }) => {
        return (
          <EditableCell
            value={row.original.status}
            type="select"
            options={processStatusOptions.map((s) => ({
              label: s,
              value: s
            }))}
            onSave={(value) =>
              updateProcessMutation.mutate({
                id: row.original.id,
                fields: { status: value }
              })
            }
          />
        )
      }
    },
    {
      header: "完成时间",
      cell: ({ row }) => {
        const startTime = row.original.startTime
        const endTime = row.original.endTime

        if (!endTime) {
          return "-"
        }

        const endDate = new Date(endTime)
        const endDateFormat = `${endDate.getFullYear()}-${
          endDate.getMonth() + 1
        }-${endDate.getDate()} ${endDate.getHours()}:${endDate.getMinutes()}`
        if (!startTime) {
          return endDateFormat
        }

        const startDate = new Date(startTime)
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
      header: "备注",
      cell: ({ row }) => {
        return (
          <EditableCell
            value={row.original.remark}
            onSave={(value) =>
              updateProcessMutation.mutate({
                id: row.original.id,
                fields: { remark: value }
              })
            }
          />
        )
      }
    }
  ]

  const table = useReactTable({
    data: data || [],
    columns,
    getCoreRowModel: getCoreRowModel()
  })

  if (isError) {
    return null
  }

  if (isPending && !data) {
    return (
      <div className="flex items-center justify-center h-full">
        <Loading content="项目流程加载中..." />
      </div>
    )
  }
  return (
    <div className="overflow-hidden rounded-md border">
      <Table>
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <TableHead key={header.id}>
                  {header.isPlaceholder
                    ? null
                    : flexRender(
                        header.column.columnDef.header,
                        header.getContext()
                      )}
                </TableHead>
              ))}
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
                无数据
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  )
}
export default ProcessPage
