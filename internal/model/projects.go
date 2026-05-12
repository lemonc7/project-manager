package model

import (
	"fmt"
	"time"

	"github.com/Oudwins/zog"
)

type CreateProjectRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Remark       string `json:"remark"`
	ManagerID    int32  `json:"managerID"`
	Priority     int32  `json:"priority"`
	DurationDays *int32 `json:"durationDays"`
}

var createProjectRequestSchema = zog.Struct(zog.Shape{
	"name":         zog.String().Required(),
	"description":  zog.String(),
	"managerID":    zog.Int32().Required(),
	"remark":       zog.String(),
	"priority":     zog.Int32().Required().GTE(1).LTE(5),
	"durationDays": zog.Ptr(zog.Int32()),
})

func (r *CreateProjectRequest) Validate() error {
	return SchemaValidate(createProjectRequestSchema, r)
}

type UpdateProjectRequest struct {
	ID           int32      `param:"id"`
	Name         *string    `json:"name"`
	Description  *string    `json:"description"`
	Remark       *string    `json:"remark"`
	StatusID     *int32     `json:"statusID"`
	ManagerID    *int32     `json:"managerID"`
	Priority     *int32     `json:"priority"`
	StartedAt    *time.Time `json:"startedAt"`
	DurationDays *int32     `json:"durationDays"`
	FinishedAt   *time.Time `json:"finishedAt"`
}

var updateProjectRequestSchema = zog.Struct(zog.Shape{
	"id":           zog.Int32().Required(),
	"name":         zog.Ptr(zog.String()),
	"description":  zog.Ptr(zog.String()),
	"remark":       zog.Ptr(zog.String()),
	"statusID":     zog.Ptr(zog.Int32()),
	"managerID":    zog.Ptr(zog.Int32()),
	"priority":     zog.Ptr(zog.Int32().GTE(1).LTE(5)),
	"startedAt":    zog.Ptr(zog.Time()),
	"durationDays": zog.Ptr(zog.Int32()),
	"finishedAt":   zog.Ptr(zog.Time()),
})

func (r *UpdateProjectRequest) Validate() error {
	return SchemaValidate(updateProjectRequestSchema, r)
}

type ListProjectsRequest struct {
	StatusIDs []int32    `json:"statusIDs"`
	ManagerID *int32     `json:"managerID"`
	StartAt   *time.Time `json:"startAt"`
	EndAt     *time.Time `json:"endAt"`
}

var listProjectRequestSchema = zog.Struct(zog.Shape{
	"statusIDs": zog.Slice(zog.Int32()),
	"managerID": zog.Ptr(zog.Int32()),
	"startAt":   zog.Ptr(zog.Time()),
	"endAt":     zog.Ptr(zog.Time()),
})

func (r *ListProjectsRequest) Validate() error {
	return SchemaValidate(listProjectRequestSchema, r)
}

type CopyProjectRequest struct {
	ID        int32  `param:"id"`
	NewName   string `json:"newName"`
	ManagerID int32  `json:"managerID"`
}

var copyProjectRequestSchema = zog.Struct(zog.Shape{
	"id":       zog.Int32().Required(),
	"newName":  zog.String().Required(),
	"mangerID": zog.Int32().Required(),
})

func (r *CopyProjectRequest) Validate() error {
	return SchemaValidate(copyProjectRequestSchema, r)
}

type IsProjectExistsRequest struct {
	NewName string `query:"newName"`
}

func (r *IsProjectExistsRequest) Validate() error {
	if r.NewName == "" {
		return fmt.Errorf("newName is empty")
	}
	return nil
}

type DeleteProjectRequest struct {
	ID int32 `param:"id"`
}

func (r *DeleteProjectRequest) Validate() error {
	return SchemaValidate(idSchema, r)
}

type StartProjectRequest struct {
	ID int32 `param:"id"`
}

func (r *StartProjectRequest) Validate() error {
	return SchemaValidate(idSchema, r)
}

type GetProjectRequest struct {
	ID int32 `param:"id"`
}

func (r *GetProjectRequest) Validate() error {
	return SchemaValidate(idSchema, r)
}

type ResetStatusRequest struct {
	ID int32 `param:"id"`
}

func (r *ResetStatusRequest) Validate() error {
	return SchemaValidate(idSchema, r)
}
