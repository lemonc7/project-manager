import { useState } from "react"
import { Input } from "../ui/input"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from "../ui/select"

type EditableType = "text" | "number" | "select"

interface EditableCellProps<T> {
  value: T
  type?: EditableType
  options?: { label: string; value: T }[]
  onSave: (value: T) => void
}

function EditableCell<T>({
  value,
  type = "text",
  options,
  onSave
}: EditableCellProps<T>) {
  const [editing, setEditing] = useState(false)
  const [tempValue, setTempValue] = useState<T>(value)

  const handleSave = () => {
    setEditing(false)
    if (tempValue !== value) {
      onSave(tempValue)
    }
  }

  // 选择框
  if (type === "select") {
    return (
      <Select
        value={String(tempValue)}
        onValueChange={(v) => {
          const newValue = v as unknown as T
          onSave(newValue)
        }}
      >
        <SelectTrigger>
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          {options?.map((opt) => (
            <SelectItem key={String(opt.value)} value={String(opt.value)}>
              {opt.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    )
  }

  if (!editing) {
    return (
      <div className="cursor-pointer" onClick={() => setEditing(true)}>
        {value !== null && value !== undefined && value !== ""
          ? String(value)
          : "-"}
      </div>
    )
  }

  // 数字输入框
  if (type === "number") {
    return (
      <Input
        type="number"
        value={tempValue as unknown as number}
        onChange={(e) => setTempValue(e.target.value as unknown as T)}
        onBlur={handleSave}
        onKeyDown={(e) => {
          if (e.key === "Enter") {
            handleSave()
          }
          if (e.key === "Escape") {
            setEditing(false)
            setTempValue(value as unknown as T)
          }
        }}
        autoFocus
      />
    )
  }

  // 文本输入框
  return (
    <Input
      value={tempValue as unknown as string}
      onChange={(e) => setTempValue(e.target.value as unknown as T)}
      onBlur={handleSave}
      onKeyDown={(e) => {
        if (e.key === "Enter") handleSave()
        if (e.key === "Escape") {
          setEditing(false)
          setTempValue(value)
        }
      }}
      autoFocus
    />
  )
}

export default EditableCell
