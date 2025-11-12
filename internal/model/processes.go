package model

import "github.com/lemonc7/project-manager/internal/repo"

type CreateProcessRequest struct {
	Name        string `json:"name" validate:"required"`
	ProjectID   int32  `json:"project_id" validate:"required"`
	EmployeeID  *int32 `json:"employee_id"`
	Description string `json:"description"`
	Sequence    int32  `json:"sequence" validate:"required"`
	Remark      string `json:"remark"`
}

type UpdateProcessRequest struct {
	ID          int32               `json:"id" validate:"required"`
	Name        *string             `json:"name"`
	Description *string             `json:"description"`
	Status      *repo.ProcessStatus `json:"status"`
	EmployeeID  *int32              `json:"employee_id"`
	Sequence    *int32              `json:"sequence"`
	Remark      *string             `json:"remark"`
}

type DeleteProcessRequest struct {
	ID int32 `json:"id" validate:"required"`
}

type ListProcessesRequest struct {
	ProjectID  *int32              `json:"project_id"`
	Status     *repo.ProcessStatus `json:"status"`
	EmployeeID *int32              `json:"employee_id"`
}
