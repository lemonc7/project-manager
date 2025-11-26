import {
  ProjectStatusOptions,
  type CreateProjectRequest,
  type Project,
  type UpdateProjectRequest
} from "@/model/projects"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogTitle
} from "../ui/dialog"
import { Field, FieldGroup, FieldLabel, FieldSet } from "../ui/field"
import { Input } from "../ui/input"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from "../ui/select"
import { Controller, useForm } from "react-hook-form"
import EmployeeSelect from "../employee/employee_select"
import useCreateProject from "@/hooks/use-create-project"
import useUpdateProject from "@/hooks/use-update-project"
import { Textarea } from "../ui/textarea"
import { Button } from "../ui/button"
import { toast } from "sonner"

type ProjectFormDialogProps =
  | {
      open: boolean
      onOpenChange: (open: boolean) => void
      mode: "create"
      defaultValues?: undefined
    }
  | {
      open: boolean
      onOpenChange: (open: boolean) => void
      mode: "edit"
      defaultValues: Project
    }

export default function ProjectFormDialog({
  open,
  onOpenChange,
  mode,
  defaultValues
}: ProjectFormDialogProps) {
  const {
    register,
    control,
    handleSubmit,
    reset,
    formState: { dirtyFields, errors }
  } = useForm<Project>({ defaultValues })

  const createProjectMutation = useCreateProject(() => {
    reset()
    onOpenChange(false)
  })
  const updateProjectMutation = useUpdateProject(() => {
    reset()
    onOpenChange(false)
  })
  const employeeID = localStorage.getItem("user_id")

  const onSubmit = (data: Project) => {
    if (mode === "create") {
      if (!employeeID) {
        toast.error("缺少user_id")
        return
      }
      const managerID = Number(employeeID)
      if (Number.isNaN(managerID) || managerID <= 0) {
        toast.error("user_id非法")
        return
      }
      const req: CreateProjectRequest = {
        name: data.name,
        description: data.description,
        managerID: managerID,
        remark: data.remark,
        duration: data.duration
      }
      createProjectMutation.mutate(req)
      return
    }
    if (!defaultValues.id) {
      console.error("缺少项目id")
      return
    }
    const req: UpdateProjectRequest = {}
    if (dirtyFields.name) req.name = data.name
    if (dirtyFields.description) req.description = data.description
    if (dirtyFields.remark) req.remark = data.remark
    if (dirtyFields.status) req.status = data.status
    if (dirtyFields.managerID) req.managerID = data.managerID
    if (dirtyFields.startedAt) req.startedAt = data.startedAt
    if (dirtyFields.duration) req.duration = data.duration
    if (dirtyFields.finishedAt) req.finishedAt = data.finishedAt

    if (Object.keys(req).length === 0) {
      onOpenChange(false)
      return
    }

    updateProjectMutation.mutate({ id: defaultValues.id, req })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogTitle>{mode === "create" ? "创建项目" : "编辑项目"}</DialogTitle>
        <DialogDescription>
          <ProjectDetails defaultValues={defaultValues} />
        </DialogDescription>
        <div className="w-full max-w-md">
          <form onSubmit={handleSubmit(onSubmit)}>
            <FieldGroup>
              <FieldSet>
                <FieldGroup>
                  {/* 项目名称 */}
                  <Field>
                    <FieldLabel htmlFor="project-name">名称</FieldLabel>
                    <Input
                      id="project-name"
                      placeholder="请输入项目名称"
                      className={errors.name ? "border-red-500" : ""}
                      {...register("name", { required: "项目名称不能为空" })}
                    />
                    {errors.name && (
                      <div className="text-red-500 text-sm mt-1">
                        {errors.name.message}
                      </div>
                    )}
                  </Field>
                  {/* 项目细节 */}
                  <Field>
                    <FieldLabel htmlFor="project-description">
                      细节描述
                    </FieldLabel>
                    <Input
                      id="project-description"
                      placeholder="细节描述"
                      {...register("description")}
                    />
                  </Field>
                  {/* 项目周期 */}
                  <Field>
                    <FieldLabel htmlFor="project-duration">项目周期</FieldLabel>
                    <Input
                      id="project-duration"
                      type="number"
                      placeholder="项目周期-天数"
                      {...register("duration", { valueAsNumber: true })}
                    />
                  </Field>
                  {/* 状态：编辑模式才显示状态 */}
                  {mode === "edit" && (
                    <div className="grid grid-cols-2 gap-4">
                      <Field>
                        <FieldLabel>状态</FieldLabel>
                        <Controller
                          control={control}
                          name="status"
                          render={({ field }) => (
                            <Select
                              value={field.value}
                              onValueChange={field.onChange}
                            >
                              <SelectTrigger>
                                <SelectValue placeholder="请选择" />
                              </SelectTrigger>
                              <SelectContent>
                                {ProjectStatusOptions.map((s) => (
                                  <SelectItem key={s} value={s}>
                                    {s}
                                  </SelectItem>
                                ))}
                              </SelectContent>
                            </Select>
                          )}
                        />
                      </Field>
                      {/* 负责人 */}
                      <Field>
                        <FieldLabel>负责人</FieldLabel>
                        <Controller
                          control={control}
                          name="managerID"
                          render={({ field }) => (
                            <EmployeeSelect
                              value={field.value}
                              onChange={field.onChange}
                            />
                          )}
                        />
                      </Field>
                    </div>
                  )}
                  {/* 备注 */}
                  <Field>
                    <FieldLabel htmlFor="project-remark">备注</FieldLabel>
                    <Textarea
                      id="project-remark"
                      placeholder="项目备注"
                      className="resize-none"
                      {...register("remark")}
                    />
                  </Field>
                </FieldGroup>
              </FieldSet>
              <Field orientation="horizontal">
                <Button type="submit">
                  {mode === "create" ? "创建" : "更新"}
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => onOpenChange(false)}
                >
                  取消
                </Button>
              </Field>
            </FieldGroup>
          </form>
        </div>
      </DialogContent>
    </Dialog>
  )
}

function ProjectDetails({ defaultValues }: { defaultValues?: Project }) {
  if (
    !defaultValues ||
    defaultValues.remainingDays === undefined ||
    !defaultValues.startedAt
  )
    return null
  const startDate = new Date(defaultValues.startedAt)
    .toISOString()
    .split("T")[0]
  return (
    <span>
      开始时间：{startDate}{" "}
      {defaultValues.remainingDays >= 0
        ? `剩余${defaultValues.remainingDays}天`
        : `超时${-defaultValues.remainingDays}天`}
    </span>
  )
}
