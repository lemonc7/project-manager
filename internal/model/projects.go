package model

import (
	"time"

	"github.com/Oudwins/zog"
	"github.com/lemonc7/project-manager/internal/repo"
)

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ManagerID   int32  `json:"managerID"`
	Remark      string `json:"remark"`
	Priority    int32  `json:"priority"`
	Duration    *int32 `json:"duration"`
}

var CreateProjectRequestSchema = zog.Struct(zog.Shape{
	"name":        zog.String().Required(),
	"description": zog.String(),
	"managerID":   zog.Int32().Required(),
	"remark":      zog.String(),
	"priority":    zog.Int32().Required().GTE(1).LTE(5),
	"duration":    zog.Int32(),
})

type UpdateProjectRequest struct {
	Name        *string             `json:"name"`
	Description *string             `json:"description"`
	Remark      *string             `json:"remark"`
	Status      *repo.ProjectStatus `json:"status"`
	ManagerID   *int32              `json:"managerID"`
	Priority    *int32              `json:"priority"`
	Duration    *int32              `json:"duration"`
}

var UpdateProjectRequestSchema = zog.Struct(zog.Shape{
	"name":        zog.String(),
	"description": zog.String(),
	"remark":      zog.String(),
	"status":      zog.String(),
	"managerID":   zog.Int32(),
	"priority":    zog.Int32().GTE(1).LTE(5),
	"duration":    zog.Int32(),
})

type ListProjectsRequest struct {
	Status    []string   `json:"status"`
	ManagerID *int32     `json:"managerID"`
	StartAt   *time.Time `json:"startAt"`
	EndAt     *time.Time `json:"endAt"`
}

type CopyProjectRequest struct {
	NewName   string `json:"newName" validate:"required"`
	ManagerID int32  `json:"managerID" validate:"required"`
}

type IsProjectExistsRequest struct {
	NewName string `query:"newName" validate:"required"`
}
