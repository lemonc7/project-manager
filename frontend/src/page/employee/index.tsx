import { deleteEmployee, listEmployees, updateEmployee } from "@/api/employees"
import { listTeams } from "@/api/teams"
import EditableCell from "@/components/common/editable-cell"
import { HeaderActions } from "@/components/common/header-actions"
import { Loading } from "@/components/common/loading"
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
import type { Employee, UpdateEmployeeRequest } from "@/model/employees"
import type { Team } from "@/model/teams"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import {
  flexRender,
  getCoreRowModel,
  useReactTable,
  type ColumnDef
} from "@tanstack/react-table"
import { MoreHorizontal } from "lucide-react"
import { useState } from "react"
import { toast } from "sonner"

export default function EmployeePage() {
  // 获取employees数据
  const {
    data: employees,
    isError,
    isLoading
  } = useQuery<Employee[], Error>({
    queryKey: ["employees", "overview"],
    queryFn: listEmployees
  })

  // 获取teams数据
  const { data: teams } = useQuery<Team[], Error>({
    queryKey: ["teams", "overview"],
    queryFn: listTeams
  })
  const queryClient = useQueryClient()
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [employeeToDelete, setEmployeeToDelete] = useState<number | null>(null)

  // 删除employee
  const deleteEmployeeMutation = useMutation({
    mutationFn: (id: number) => deleteEmployee(id),
    onSuccess: (_, deleteEmployeeID) => {
      queryClient.setQueryData<Employee[]>(
        ["employees", "overview"],
        (currentEmployees) => {
          if (!currentEmployees) {
            return currentEmployees
          }
          return currentEmployees.filter(
            (employee) => employee.id !== deleteEmployeeID
          )
        }
      )
      toast.success("删除成功")
      setIsDialogOpen(false)
      setEmployeeToDelete(null)
    },
    onError: () => {
      toast.error("删除失败")
      setIsDialogOpen(false)
      setEmployeeToDelete(null)
    }
  })

  // 更新employee
  const updateEmployeeMutation = useMutation({
    mutationFn: ({
      id,
      fields
    }: {
      id: number
      fields: Partial<UpdateEmployeeRequest>
    }) => updateEmployee(id, fields),
    onSuccess: (updatedEmployee) => {
      // 更新缓存，只更新当前成员的数据
      queryClient.setQueryData<Employee[]>(
        ["employees", "overview"],
        (currentEmployees) =>
          currentEmployees
            ? currentEmployees.map((employee) =>
                employee.id === updatedEmployee.id
                  ? { ...employee, ...updatedEmployee }
                  : employee
              )
            : [updatedEmployee]
      )
    },
    onError: () => {
      toast.error("更新失败")
    }
  })

  const handleDeleteConfirm = () => {
    if (employeeToDelete !== null) {
      deleteEmployeeMutation.mutate(employeeToDelete)
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
        if (!teams) {
          return <span>{row.original.teamId}</span>
        }
        return (
          <EditableCell
            value={row.original.teamId ?? undefined}
            type="select"
            options={teams.map((t) => ({
              label: t.name,
              value: t.id
            }))}
            onSave={(value) =>
              updateEmployeeMutation.mutate({
                id: row.original.id,
                fields: { teamID: Number(value) }
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
                  setEmployeeToDelete(employee.id)
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
    return null
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
            <AlertDialogTitle>确定要删除该成员吗?</AlertDialogTitle>
            <AlertDialogDescription>
              这个操作不能被撤销，将从数据库中删除该成员，你确定要继续吗？
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel
              onClick={() => {
                setEmployeeToDelete(null)
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
