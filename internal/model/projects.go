package model

import (
	"time"

	"github.com/lemonc7/project-manager/internal/repo"
)

type CreateProjectRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
	ManagerID   int32   `json:"managerID" validate:"required"`
	Remark      *string `json:"remark"`
	Duration    *int32  `json:"duration"`
}

type UpdateProjectRequest struct {
	Name        *string             `json:"name"`
	Description *string             `json:"description"`
	Remark      *string             `json:"remark"`
	Status      *repo.ProjectStatus `json:"status"`
	ManagerID   *int32              `json:"managerID"`
	StartedAt   *time.Time          `json:"startedAt"`
	Duration    *int32              `json:"duration"`
	FinishedAt  *time.Time          `json:"finishedAt"`
}

type ListProjectsRequest struct {
	Status    *repo.ProjectStatus `query:"status"`
	ManagerID *int32              `query:"managerID"`
	StartAt   *time.Time          `query:"startAt"`
	EndAt     *time.Time          `query:"endAt"`
}
