import useTeams from "@/hooks/use-teams"
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue
} from "../ui/select"

export default function TeamSelect({
  value,
  onChange
}: {
  value?: string
  onChange?: (value: string) => void
}) {
  const { data, isLoading, isError } = useTeams()

  if (isLoading) {
    return <div>加载中</div>
  }
  if (isError) {
    return <div>加载失败</div>
  }

  return (
    <Select value={value} onValueChange={onChange}>
      <SelectTrigger>
        <SelectValue placeholder="请选择" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          <SelectLabel>团队列表</SelectLabel>
          {data?.map((team) => (
            <SelectItem key={team.id} value={String(team.id)}>
              {team.name}
            </SelectItem>
          ))}
        </SelectGroup>
      </SelectContent>
    </Select>
  )
}
