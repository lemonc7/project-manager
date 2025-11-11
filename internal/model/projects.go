package model

import (
	"time"

	"github.com/lemonc7/project-manager/internal/repo"
)

type CreateProjectRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
	ManagerID   *int32 `json:"manager_id"`
}

type UpdateProjectRequest struct {
	ID          int32              `json:"id" validate:"required"`
	Name        string             `json:"name" validate:"required"`
	Description string             `json:"description,omitempty"`
	Status      repo.ProcessStatus `json:"status"`
	ManagerID   *int32             `json:"manager_id,omitempty"`
	StartTime   *time.Time         `json:"start_time"`
	Duration    int32              `json:"duration"`
}

type DeleteProjectRequest struct {
	ID int32 `json:"id" validate:"required"`
}


