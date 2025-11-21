package model

import (
	"time"

	"github.com/lemonc7/project-manager/internal/repo"
)

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
	ProjectID  *int32              `query:"projectID"`
	Status     *repo.ProcessStatus `query:"status"`
	EmployeeID *int32              `query:"employeeID"`
}

type Task struct {
	ID            int32      `json:"id"`
	Name          string     `json:"name"`
	StartTime     *time.Time `json:"startTime,omitempty"`
	RemainingDays *int32     `json:"remainingDays,omitempty"`
	ProjectID     int32      `json:"projectID"`
	ProjectName   string     `json:"projectName"`
}
