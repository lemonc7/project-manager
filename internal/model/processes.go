package model

import (
	"github.com/lemonc7/project-manager/internal/repo"
)

type InsertProcessRequest struct {
	Name        string `json:"name" validate:"required"`
	ProjectID   int32  `json:"projectID" validate:"required"`
	EmployeeID  int32  `json:"employeeID"`
	Description string `json:"description"`
	Remark      string `json:"remark"`
}

type UpdateProcessRequest struct {
	Name        *string             `json:"name"`
	Description *string             `json:"description"`
	Remark      *string             `json:"remark"`
	Status      *repo.ProcessStatus `json:"status"`
	EmployeeID  *int32              `json:"employeeID"`
	SpendTime   *float64            `json:"spendTime" validate:"gte=0"`
}

type FinishProcessRequest struct {
	ProjectID int32   `json:"projectID" validate:"required"`
	Sequence  int32   `json:"sequence" validate:"required"`
	SpendTime float64 `json:"spendTime" validate:"gte=0"`
	Remark    *string `json:"remark"`
}

type ReworkProcessRequest struct {
	SpendTime float64 `json:"spendTime" validate:"gte=0"`
	Remark    string  `json:"remark"`
}

type DeleteProcessesRequest struct {
	Sequences []int32 `json:"sequences" validate:"required"`
}
