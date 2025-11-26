import EditableCell from "@/components/common/editable-cell"
import HeaderActions from "@/components/common/header-actions"
import Loading from "@/components/common/loading"
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
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from "@/components/ui/dropdown-menu"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from "@/components/ui/table"
import useDeleteEmployee from "@/hooks/use-delete-employee"
import useEmployees from "@/hooks/use-employees"
import useTeams from "@/hooks/use-teams"
import useUpdateEmployee from "@/hooks/use-update-employee"
import type { Employee } from "@/model/employees"
import {
  flexRender,
  getCoreRowModel,
  useReactTable,
  type ColumnDef
} from "@tanstack/react-table"
import { MoreHorizontal } from "lucide-react"
import { useState } from "react"

export default function EmployeePage() {
  // 获取employees数据
  const { data: employees, isError, isLoading } = useEmployees()

  // 获取teams数据
  const { data: teams } = useTeams()

  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [selectedEmployee, setSelectedEmployee] = useState<Employee | null>(
    null
  )

  // 删除employee
  const deleteEmployeeMutation = useDeleteEmployee(
    () => {
      setIsDialogOpen(false)
      setSelectedEmployee(null)
    },
    () => {
      setIsDialogOpen(false)
      setSelectedEmployee(null)
    }
  )

  // 更新employee
  const updateEmployeeMutation = useUpdateEmployee()
  const handleDeleteConfirm = () => {
    if (selectedEmployee !== null) {
      deleteEmployeeMutation.mutate(selectedEmployee.id)
    }
  }

  // table列定义
  const columns: ColumnDef<Employee>[] = [
    {
      accessorKey: "name",
      header: "成员"
    },
    {
      header: "小组",
      cell: ({ row }) => {
        return (
          <EditableCell
            value={row.original.teamID}
            type="select"
            options={teams?.map((t) => ({
              label: t.name,
              value: t.id
            }))}
            onSave={(value) =>
              updateEmployeeMutation.mutate({
                id: row.original.id,
                req: { teamID: value }
              })
            }
          />
        )
      }
    },
    {
      id: "actions",
      header: "操作",
      size: 30,
      cell: ({ row }) => {
        const employee = row.original
        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <MoreHorizontal />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuLabel>操作</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                onClick={() => console.log("点击了编辑", employee.id)}
              >
                编辑
              </DropdownMenuItem>
              <DropdownMenuItem
                className="text-red-600 focus:text-red-600"
                onClick={() => {
                  setSelectedEmployee(employee)
                  setIsDialogOpen(true)
                }}
              >
                删除
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        )
      }
    }
  ]

  const table = useReactTable({
    data: employees || [],
    columns,
    getCoreRowModel: getCoreRowModel()
  })

  if (isError) {
    return (
      <div className="flex items-center justify-center h-full">
        <span>加载失败</span>
      </div>
    )
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-full">
        <Loading content="成员列表加载中..." />
      </div>
    )
  }

  return (
    <>
      <HeaderActions>
        <Button>新建</Button>
      </HeaderActions>
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
      {/* 弹窗 */}
      <AlertDialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
        <AlertDialogContent onOpenAutoFocus={(e) => e.preventDefault()}>
          <AlertDialogHeader>
            <AlertDialogTitle>删除成员</AlertDialogTitle>
            <AlertDialogDescription>
              确定要删除{`[${selectedEmployee?.name}]`}吗?
              这个操作不能被撤销，将从数据库中删除该成员。
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel
              onClick={() => {
                setSelectedEmployee(null)
              }}
              disabled={deleteEmployeeMutation.isPending}
            >
              取消
            </AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDeleteConfirm}
              disabled={deleteEmployeeMutation.isPending}
            >
              {deleteEmployeeMutation.isPending ? "删除中..." : "确认"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  )
}
