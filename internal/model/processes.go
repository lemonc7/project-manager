package model

import "github.com/lemonc7/project-manager/internal/repo"

type InsertProcessRequest struct {
	Name        string `json:"name" validate:"required"`
	ProjectID   int32  `json:"projectID" validate:"required"`
	EmployeeID  int32  `json:"employeeID" validate:"required"`
	Description string `json:"description"`
	Remark      string `json:"remark"`
}

type UpdateProcessRequest struct {
	Name        *string             `json:"name"`
	Description *string             `json:"description"`
	Status      *repo.ProcessStatus `json:"status"`
	EmployeeID  *int32              `json:"employeeID"`
	Remark      *string             `json:"remark"`
}

type ListProcessesRequest struct {
	ProjectID  *int32              `json:"projectID"`
	Status     *repo.ProcessStatus `json:"status"`
	EmployeeID *int32              `json:"employeeID"`
}
