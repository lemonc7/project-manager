import { ProjectStatusOptions, type Project } from "@/model/projects"
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
import EditableCell from "@/components/common/editable-cell"
import Loading from "@/components/common/loading"
import HeaderActions from "@/components/common/header-actions"
import { useState } from "react"
import { Button } from "@/components/ui/button"
import ProjectFormDialog from "@/components/project/project-form-dialog"
import useUpdateProject from "@/hooks/use-update-project"
import useProjects from "@/hooks/use-projects"
import useEmployees from "@/hooks/use-employees"

export default function ProjectPage() {
  const { data, isError, isPending } = useProjects({})
  const [createDialog, setCreateDialog] = useState(false)

  const updateProjectMutation = useUpdateProject()
  const { data: employees } = useEmployees()

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
              req: { description: value }
            })
          }
        />
      )
    },
    {
      header: "负责人",
      cell: ({ row }) => (
        <EditableCell
          value={row.original.managerID}
          type="select"
          options={employees?.map((emp) => ({
            label: emp.name,
            value: emp.id
          }))}
          disabled
        />
      )
    },
    {
      header: "开始时间",
      cell: ({ row }) => {
        let value = ""
        if (row.original.startedAt) {
          value = new Date(row.original.startedAt).toISOString().split("T")[0]
        }
        if (row.original.remainingDays) {
          if (row.original.remainingDays < 0) {
            value += ` (超时${-row.original.remainingDays}天)`
          } else {
            value += ` (剩余${row.original.remainingDays}天)`
          }
        }
        return <span>{value}</span>
      }
    },
    {
      header: "状态",
      cell: ({ row }) => (
        <EditableCell
          value={row.original.status}
          type="select"
          options={ProjectStatusOptions.map((s) => ({
            label: s,
            value: s
          }))}
          onSave={(value) =>
            updateProjectMutation.mutate({
              id: row.original.id,
              req: { status: value }
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
              req: { remark: value }
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

  if (isPending && !data) {
    return (
      <div className="flex items-center justify-center h-full">
        <Loading content="项目列表加载中..." />
      </div>
    )
  }

  if (isError) {
    return (
      <div className="flex items-center justify-center h-full">
        <span>加载失败</span>
      </div>
    )
  }
  return (
    <>
      <HeaderActions>
        <Button onClick={() => setCreateDialog(true)}>新建项目</Button>
      </HeaderActions>
      <ProjectFormDialog
        open={createDialog}
        onOpenChange={setCreateDialog}
        mode="create"
      />
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
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext()
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="h-24 text-center"
                >
                  无数据
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
    </>
  )
}
