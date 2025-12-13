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
	Priority    int32  `json:"priority" validate:"required,gte=1,lte=5"`
	Duration    *int32 `json:"duration"`
}

type UpdateProjectRequest struct {
	Name        *string             `json:"name"`
	Description *string             `json:"description"`
	Remark      *string             `json:"remark"`
	Status      *repo.ProjectStatus `json:"status"`
	ManagerID   *int32              `json:"managerID"`
	Priority    *int32              `json:"priority"`
	Duration    *int32              `json:"duration"`
}

type ListProjectsRequest struct {
	StartAt *time.Time `query:"startAt"`
	EndAt   *time.Time `query:"endAt"`
}

type CopyProjectRequest struct {
	NewName   string `json:"newName" validate:"required"`
	ManagerID int32  `json:"managerID" validate:"required"`
}

type IsProjectExistsRequest struct {
	NewName string `query:"newName" validate:"required"`
}
