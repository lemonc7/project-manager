import useEmployees from "@/hooks/use-employees"
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue
} from "../ui/select"

export default function EmployeeSelect({
  value,
  onChange,
  disabled = false
}: {
  value?: number
  onChange?: (value: number) => void
  disabled?: boolean
}) {
  const { data, isLoading, isError } = useEmployees()

  if (isLoading) {
    return <div>加载中</div>
  }
  if (isError) {
    return <div>加载失败</div>
  }

  return (
    <Select
      value={value ? String(value) : undefined}
      onValueChange={(v) => onChange?.(Number(v))}
      disabled={disabled}
    >
      <SelectTrigger>
        <SelectValue placeholder="请选择" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          <SelectLabel>成员列表</SelectLabel>
          {data?.map((emp) => (
            <SelectItem key={emp.id} value={String(emp.id)}>
              {emp.name}
            </SelectItem>
          ))}
        </SelectGroup>
      </SelectContent>
    </Select>
  )
}
