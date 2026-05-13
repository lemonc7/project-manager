package model

import (
	"github.com/lemonc7/project-manager/internal/repo"
)

type InsertProcessRequest struct {
	StepOrder   int32  `param:"stepOrder" validate:"required,gte=1"`
	Name        string `json:"name" validate:"required"`
	ProjectID   int32  `json:"projectID" validate:"required"`
	EmployeeID  int32  `json:"employeeID" validate:"required"`
	Description string `json:"description"`
	Remark      string `json:"remark"`
}

type UpdateProcessRequest struct {
	ID          int32    `param:"id" validate:"required"`
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Remark      *string  `json:"remark"`
	StatusID    *int32   `json:"statusID"`
	EmployeeID  *int32   `json:"employeeID"`
	SpendHours  *float64 `json:"spendHours"`
}

type FinishProcessRequest struct {
	ID         int32   `param:"id" validate:"required"`
	ProjectID  int32   `json:"projectID" validate:"required"`
	StepOrder  int32   `json:"stepOrder" validate:"required,gte=1"`
	SpendHours float64 `json:"spendHours" validate:"gte=0"`
	Remark     *string `json:"remark"`
}

type ReworkProcessRequest struct {
	ID         int32   `param:"id" validate:"required"`
	SpendHours float64 `json:"spendHours" validate:"gte=0"`
	Remark     string  `json:"remark"`
}

type DeleteProcessRequest struct {
	ID int32 `param:"id" validate:"required"`
}

type DeleteProcessesRequest struct {
	ProjectID  int32   `param:"projectID" validate:"required"`
	StepOrders []int32 `json:"stepOrders" validate:"required,dive,required,gte=1"`
}

type ListProcessesRequest struct {
	ID int32 `param:"id" validate:"required"`
}

type ListTasksRequest struct {
	EmployeeID int32               `param:"id" validate:"required"`
	Status     repo.StatusCategory `query:"status" validate:"oneof=todo ready doing finished"`
}
