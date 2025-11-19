import { listProjects, updateProject } from "@/api/projects"
import type {
  Project,
  ProjectStatus,
  UpdateProjectRequest
} from "@/model/projects"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
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
import EditableCell from "@/components/editable-cell"

const projectStatusOptions = [
  "排单",
  "进行",
  "暂停",
  "完成",
  "终止",
  "归档"
] satisfies ProjectStatus[]

function ProjectPage() {
  const { data, isError, isPending } = useQuery<Project[], Error>({
    queryKey: ["projects", "overview"],
    queryFn: listProjects
  })

  // const {
  //   data: employees = [],
  //   isError: isEmployeeError,
  //   isPending: isEmployeePending
  // } = useQuery<Employee[], Error>({
  //   queryKey: ["employees"],
  //   queryFn: listEmployees
  // })

  const queryClient = useQueryClient()

  const updateProjectMutation = useMutation({
    mutationFn: ({
      id,
      fields
    }: {
      id: number
      fields: Partial<Omit<UpdateProjectRequest, "id">>
    }) => updateProject({ id, ...fields }),
    onSuccess: (updatedProject) => {
      // 更新缓存: 只更新当前项目的数据
      queryClient.setQueryData<Project[]>(
        ["projects", "overview"],
        (currentProjects) =>
          currentProjects
            ? currentProjects.map((project) =>
                project.id === updatedProject.id
                  ? { ...project, ...updatedProject }
                  : project
              )
            : [updatedProject]
      )
    },
    onError: () => {
      toast.error("更新失败")
    }
  })
  const navigate = useNavigate()

  const columns: ColumnDef<Project>[] = [
    {
      header: "项目",
      cell: ({ row }) => (
        <div
          className="cursor-pointer"
          onClick={() => navigate(`/projects/${row.original.id}`)}
        >
          {row.original.name}
        </div>
      )
    },
    {
      header: "描述",
      cell: ({ row }) => (
        <EditableCell
          value={row.original.description}
          onSave={(value) =>
            updateProjectMutation.mutate({
              id: row.original.id,
              fields: { description: value }
            })
          }
        />
      )
    },
    {
      accessorKey: "employeeName",
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
      header: "状态",
      cell: ({ row }) => (
        <EditableCell
          value={row.original.status}
          type="select"
          options={projectStatusOptions.map((s) => ({
            label: s,
            value: s
          }))}
          onSave={(value) =>
            updateProjectMutation.mutate({
              id: row.original.id,
              fields: { status: value }
            })
          }
        />
      )
    },
    {
      header: "备注",
      cell: ({ row }) => (
        <EditableCell
          value={row.original.remark}
          onSave={(value) =>
            updateProjectMutation.mutate({
              id: row.original.id,
              fields: { remark: value }
            })
          }
        />
      )
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
export default ProjectPage
