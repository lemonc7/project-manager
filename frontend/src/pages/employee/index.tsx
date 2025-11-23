import { deleteEmployee, listEmployees } from "@/api/employees"
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
import type { Employee } from "@/models/employees"
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

function EmployeePage() {
  // 获取employees数据
  const { data, isError, isLoading } = useQuery<Employee[], Error>({
    queryKey: ["employees"],
    queryFn: () => {
      const api = listEmployees()
      toast.promise(api, {
        loading: "成员加载中...",
        success: "成员加载成功",
        error: "加载失败"
      })
      return api
    },
    select: (data) => data.filter((employee) => employee.isActive === true)
  })

  const queryClient = useQueryClient()
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [employeeToDelete, setEmployeeToDelete] = useState<number | null>(null)

  // 删除employee函数
  const deleteEmployeeMutation = useMutation({
    mutationFn: (id: number) => deleteEmployee(id),
    onSuccess: (_, deleteEmployeeID) => {
      queryClient.setQueryData<Employee[]>(
        ["employees"],
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

  const handleDeleteConfirm = () => {
    if (employeeToDelete !== null) {
      deleteEmployeeMutation.mutate(employeeToDelete)
    }
  }

  // table列定义
  const columns: ColumnDef<Employee>[] = [
    {
      accessorKey: "employeeName",
      header: "成员"
    },
    {
      accessorKey: "teamName",
      header: "小组"
    },
    {
      header: "添加时间",
      cell: ({ row }) => {
        const createdAt = row.original.createdAt
        return createdAt ? formatDate(createdAt) : "-"
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
    data: data || [],
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

function formatDate(dateStr: string) {
  const date = new Date(dateStr)
  const year = date.getFullYear()
  const month = date.getMonth() + 1
  const day = date.getDate()
  return `${year}-${month}-${day}`
}

export default EmployeePage
