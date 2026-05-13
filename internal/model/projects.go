package model

import (
	"time"
)

type CreateProjectRequest struct {
	Name         string `json:"name" validate:"required"`
	Description  string `json:"description"`
	Remark       string `json:"remark"`
	ManagerID    int32  `json:"managerID" validate:"required"`
	Priority     int32  `json:"priority" validate:"required,gte=1,lte=5"`
	DurationDays *int32 `json:"durationDays"`
}

type UpdateProjectRequest struct {
	ID           int32      `param:"id" validate:"required"`
	Name         *string    `json:"name"`
	Description  *string    `json:"description"`
	Remark       *string    `json:"remark"`
	StatusID     *int32     `json:"statusID"`
	ManagerID    *int32     `json:"managerID"`
	Priority     *int32     `json:"priority" validate:"gte=1,lte=5"`
	StartedAt    *time.Time `json:"startedAt"`
	DurationDays *int32     `json:"durationDays"`
	FinishedAt   *time.Time `json:"finishedAt"`
}

type ListProjectsRequest struct {
	StatusIDs []int32    `json:"statusIDs" validate:"required"`
	ManagerID *int32     `json:"managerID"`
	StartAt   *time.Time `json:"startAt"`
	EndAt     *time.Time `json:"endAt"`
}

type CopyProjectRequest struct {
	ID        int32  `param:"id" validate:"required"`
	NewName   string `json:"newName" validate:"required"`
	ManagerID int32  `json:"managerID" validate:"required"`
}

type IsProjectExistsRequest struct {
	NewName string `query:"newName" validate:"required"`
}

type DeleteProjectRequest struct {
	ID int32 `param:"id" validate:"required"`
}

type StartProjectRequest struct {
	ID int32 `param:"id" validate:"required"`
}

type GetProjectRequest struct {
	ID int32 `param:"id" validate:"required"`
}

type ResetStatusRequest struct {
	ID int32 `param:"id" validate:"required"`
}
