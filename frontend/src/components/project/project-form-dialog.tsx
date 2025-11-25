import {
  ProjectStatusOptions,
  type CreateProjectRequest,
  type Project,
  type UpdateProjectRequest
} from "@/model/projects"
import { Dialog, DialogContent } from "../ui/dialog"
import {
  Field,
  FieldDescription,
  FieldGroup,
  FieldLabel,
  FieldLegend,
  FieldSet
} from "../ui/field"
import { Input } from "../ui/input"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from "../ui/select"
import { Controller, useForm } from "react-hook-form"

type ProjectFormDialogProps =
  | {
      open: boolean
      onOpenChange: (open: boolean) => void
      mode: "create"
      defaultValues?: undefined
      submit: (data: CreateProjectRequest) => void
    }
  | {
      open: boolean
      onOpenChange: (open: boolean) => void
      mode: "edit"
      defaultValues: Project
      submit: (data: UpdateProjectRequest) => void
    }

export default function ProjectFormDialog({
  open,
  onOpenChange,
  mode,
  defaultValues,
  submit
}: ProjectFormDialogProps) {
  const {
    register,
    control,
    handleSubmit,
    formState: { dirtyFields }
  } = useForm<Project>({ defaultValues })

  const onSubmit = (data: Project) => {
    if (mode === "create") {
      const req: CreateProjectRequest = {
        name: data.name,
        description: data.description,
        managerID: data.managerID,
        remark: data.remark,
        duration: data.duration
      }
      submit(req)
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
    submit(req)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <div className="w-full max-w-md">
          <form onSubmit={handleSubmit(onSubmit)}>
            <FieldGroup>
              <FieldSet>
                <FieldLegend>
                  {mode === "create" ? "创建项目" : "编辑项目"}
                </FieldLegend>
                {defaultValues && <FieldDescription></FieldDescription>}
                <FieldGroup>
                  {/* 项目名称 */}
                  <Field>
                    <FieldLabel htmlFor="project-name">名称</FieldLabel>
                    <Input
                      id="project-name"
                      placeholder="请输入项目名称"
                      {...register("name", { required: true })}
                    />
                  </Field>
                  {/* 项目细节 */}
                  <Field>
                    <FieldLabel htmlFor="project-description">
                      细节描述
                    </FieldLabel>
                    <Input
                      id="project-description"
                      placeholder="项目细节描述"
                      {...register("description")}
                    />
                  </Field>
                  {/* 状态：编辑模式才显示状态 */}
                  {mode === "edit" && (
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
                  )}
                  <Field>
                    <FieldLabel>负责人</FieldLabel>
                  </Field>
                </FieldGroup>
              </FieldSet>
            </FieldGroup>
          </form>
        </div>
      </DialogContent>
    </Dialog>
  )
}
