import {
  ProjectStatusOptions,
  type CreateProjectRequest,
  type Project,
  type ProjectStatus,
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
import { useForm } from "react-hook-form"

interface ProjectFormValues {
  name: string
  description?: string
  remark?: string
  status: ProjectStatus
  managerID: number
  startedAt?: string
  duration?: number
  finishedAt?: string
}

function mapProjectToFormValues(project: Project): ProjectFormValues {
  return {
    name: project.name,
    description: project.description ?? "",
    remark: project.remark ?? "",
    status: project.status,
    managerID: project.managerId,
    startedAt: project.startedAt ?? "",
    duration: project.duration ?? undefined,
    finishedAt: project.finishedAt ?? ""
  }
}

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
    handleSubmit,
    formState: { dirtyFields }
  } = useForm({
    defaultValues: defaultValues
      ? mapProjectToFormValues(defaultValues)
      : undefined
  })

  const onSubmit = (data: ProjectFormValues) => {
    if (mode === "create") {
      const req: CreateProjectRequest = {
        name: data.name,
        description: data.description,
        managerID: data.managerID,
        remark: data.remark,
        duration: data.duration
      }
      submit(req)
    } else {
      const req: UpdateProjectRequest = {}

      submit(req)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <div className="w-full max-w-md">
          <form>
            <FieldGroup>
              <FieldSet>
                <FieldLegend>
                  {mode === "create" ? "创建项目" : "编辑项目"}
                </FieldLegend>
                {defaultValues && <FieldDescription></FieldDescription>}

                <FieldGroup>
                  <Field>
                    <FieldLabel htmlFor="project-name">名称</FieldLabel>
                    <Input
                      id="project-name"
                      name="name"
                      placeholder="请输入项目名称"
                      value={defaultValues?.name}
                      required
                    />
                  </Field>
                  <Field>
                    <FieldLabel htmlFor="project-description">
                      细节描述
                    </FieldLabel>
                    <Input
                      id="project-description"
                      name="description"
                      placeholder="项目细节描述"
                      value={defaultValues?.description ?? undefined}
                    />
                  </Field>
                  <Field>
                    <FieldLabel>状态</FieldLabel>
                    <Select
                      name="status"
                      value={mode === "edit" ? defaultValues?.status : "排单"}
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
