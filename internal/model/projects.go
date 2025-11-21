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
}

type UpdateProjectRequest struct {
	Name        *string             `json:"name"`
	Description *string             `json:"description"`
	Status      *repo.ProjectStatus `json:"status"`
	ManagerID   *int32              `json:"managerID"`
	StartTime   *time.Time          `json:"startTime"`
	Duration    *int32              `json:"duration"`
	Remark      *string             `json:"remark"`
}

type ListProjectsRequest struct {
	Status    *repo.ProjectStatus `query:"status"`
	ManagerID *int32              `query:"managerID"`
	StartTime *time.Time          `query:"startTime"`
	EndTime   *time.Time          `query:"endTime"`
}
