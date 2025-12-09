package model

import (
	"time"

	"github.com/lemonc7/project-manager/internal/repo"
)

type CreateProjectRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	ManagerID   int32  `json:"managerID" validate:"required"`
	Remark      string `json:"remark"`
	Priority    int32  `json:"priority" validate:"required,min=1,max=5"`
	Duration    *int32 `json:"duration"`
}

type UpdateProjectRequest struct {
	Name        *string             `json:"name"`
	Description *string             `json:"description"`
	Remark      *string             `json:"remark"`
	Status      *repo.ProjectStatus `json:"status"`
	ManagerID   *int32              `json:"managerID"`
	Priority    *int32              `json:"priority"`
	StartedAt   *time.Time          `json:"startedAt"`
	Duration    *int32              `json:"duration"`
	FinishedAt  *time.Time          `json:"finishedAt"`
}

type ListProjectsRequest struct {
	StartAt *time.Time `query:"startAt"`
	EndAt   *time.Time `query:"endAt"`
}
