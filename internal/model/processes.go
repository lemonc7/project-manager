package model

import (
	"time"

	"github.com/lemonc7/project-manager/internal/repo"
)

type InsertProcessRequest struct {
	Name        string  `json:"name" validate:"required"`
	ProjectID   int32   `json:"projectID" validate:"required"`
	EmployeeID  *int32  `json:"employeeID"`
	Description *string `json:"description"`
	Remark      *string `json:"remark"`
}

type UpdateProcessRequest struct {
	Name        *string             `json:"name"`
	Description *string             `json:"description"`
	Remark      *string             `json:"remark"`
	Status      *repo.ProcessStatus `json:"status"`
	EmployeeID  *int32              `json:"employeeID"`
	SpendTime   *float64            `json:"spendTime"`
	StartedAt   *time.Time          `json:"startedAt"`
	FinishedAt  *time.Time          `json:"finishedAt"`
}

type FinishProcessRequest struct {
	SpendTime float64 `json:"spendTime" validate:"required"`
}

type Process struct {
	ID          int32              `json:"id"`
	ProjectID   int32              `json:"projectID"`
	Name        string             `json:"name"`
	Description *string            `json:"description,omitempty"`
	Remark      *string            `json:"remark,omitempty"`
	Status      repo.ProcessStatus `json:"status"`
	EmployeeID  *int32             `json:"employeeID,omitempty"`
	Sequence    int32              `json:"sequence"`
	SpendTime   *float64           `json:"spendTime,omitempty"`
	StartedAt   *time.Time         `json:"startedAt,omitempty"`
	FinishedAt  *time.Time         `json:"finishedAt,omitempty"`
}

type Task struct {
	ID            int32      `json:"id"`
	Name          string     `json:"name"`
	Description   *string    `json:"description,omitempty"`
	StartedAt     *time.Time `json:"startedAt,omitempty"`
	RemainingDays *int32     `json:"remainingDays,omitempty"`
	ProjectID     int32      `json:"projectID"`
	ProjectName   string     `json:"projectName"`
}

func (t Task) GetRemainingDays() *int32 {
	return t.RemainingDays
}

type HasRemainingDays interface {
	GetRemainingDays() *int32
}
