package model

type ListTodosRequest struct {
	Status []string `json:"status"`
}

type CreateTodoRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Remark      string `json:"remark"`
	Priority    int32  `json:"priority" validate:"required,gte=1,lte=5"`
	ReporterID  int32  `json:"reporterID" validate:"required"`
	AssigneeID  int32  `json:"assigneeID" validate:"required"`
}

type UpdateTodoRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Remark      *string `json:"remark"`
	Priority    *int32  `json:"priority"`
	Status      *string `json:"status"`
	ReporterID  *int32  `json:"reporterID"`
	AssigneeID  *int32  `json:"assigneeID"`
}
