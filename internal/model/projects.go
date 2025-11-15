package model

import (
	"time"

	"github.com/lemonc7/project-manager/internal/repo"
)

type CreateProjectRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	ManagerID   int32  `json:"manager_id" validate:"required"`
	Remark      string `json:"remark"`
}

type UpdateProjectRequest struct {
	Name        *string             `json:"name"`
	Description *string             `json:"description"`
	Status      *repo.ProcessStatus `json:"status"`
	ManagerID   *int32              `json:"manager_id"`
	StartTime   *time.Time          `json:"start_time"`
	Duration    *int32              `json:"duration"`
	Remark      *string             `json:"remark"`
}

type ListProjectsRequest struct {
	Status    *repo.ProjectStatus `json:"status"`
	ManagerID *int32              `json:"manager_id"`
	StartTime *time.Time          `json:"start_time"`
	EndTime   *time.Time          `json:"end_time"`
}
