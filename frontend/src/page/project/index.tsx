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
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from "@/components/ui/dropdown-menu"
import { MoreHorizontal } from "lucide-react"
import useDeleteProject from "@/hooks/use-delete-project"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle
} from "@/components/ui/alert-dialog"

export default function ProjectPage() {
  const { data, isError, isPending } = useProjects({})
  const { data: employees } = useEmployees()

  const [createDialog, setCreateDialog] = useState(false)
  const [editDialog, setEditDialog] = useState(false)
  const [deleteDialog, setDeleteDialog] = useState(false)
  const [selectedProject, setSelectedProject] = useState<Project | null>(null)

  const updateProjectMutation = useUpdateProject(
    () => {
      setEditDialog(false)
      setSelectedProject(null)
    },
    () => {
      setEditDialog(false)
      setSelectedProject(null)
    }
  )
  const deleteProjectMutation = useDeleteProject(
    () => {
      setDeleteDialog(false)
      setSelectedProject(null)
    },
    () => {
      setDeleteDialog(false)
      setSelectedProject(null)
    }
  )
  const handleDeleteConfirm = () => {
    if (selectedProject !== null) {
      deleteProjectMutation.mutate(selectedProject.id)
    }
  }
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
          disabled
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
    },
    {
      id: "actions",
      header: "操作",
      cell: ({ row }) => {
        return (
          <>
            <div>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" className="h-8 w-8 p-0">
                    <MoreHorizontal />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent>
                  <DropdownMenuItem>开始</DropdownMenuItem>
                  <DropdownMenuItem
                    onClick={() => {
                      setSelectedProject(row.original)
                      setEditDialog(true)
                    }}
                  >
                    编辑
                  </DropdownMenuItem>
                  <DropdownMenuItem
                    className="text-red-600 focus:text-red-600"
                    onClick={() => {
                      setSelectedProject(row.original)
                      setDeleteDialog(true)
                    }}
                  >
                    删除
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </>
        )
      }
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
      {/* 编辑弹窗 */}
      {selectedProject && (
        <ProjectFormDialog
          open={editDialog}
          onOpenChange={setEditDialog}
          mode="edit"
          defaultValues={selectedProject}
        />
      )}
      {/* 创建弹窗 */}
      <ProjectFormDialog
        open={createDialog}
        onOpenChange={setCreateDialog}
        mode="create"
      />
      {/* 删除弹窗 */}
      <AlertDialog open={deleteDialog} onOpenChange={setDeleteDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>删除项目</AlertDialogTitle>
            <AlertDialogDescription>
              确定要删除项目
              <span className="font-bold text-red-600">
                {selectedProject?.name}
              </span>
              吗? 这个操作不能被撤销，将从数据库中删除该项目以及
              <span className="font-bold text-red-600">相关流程</span>。
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel
              onClick={() => {
                setSelectedProject(null)
              }}
              disabled={deleteProjectMutation.isPending}
            >
              取消
            </AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDeleteConfirm}
              disabled={deleteProjectMutation.isPending}
            >
              {deleteProjectMutation.isPending ? "删除中..." : "确认"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
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
