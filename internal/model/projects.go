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
	"Name":         zog.String().Required(),
	"Description":  zog.String(),
	"ManagerID":    zog.Int32().Required(),
	"Remark":       zog.String(),
	"Priority":     zog.Int32().Required().GTE(1).LTE(5),
	"DurationDays": zog.Ptr(zog.Int32()),
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
	"ID":           zog.Int32().Required(),
	"Name":         zog.Ptr(zog.String()),
	"Description":  zog.Ptr(zog.String()),
	"Remark":       zog.Ptr(zog.String()),
	"StatusID":     zog.Ptr(zog.Int32()),
	"ManagerID":    zog.Ptr(zog.Int32()),
	"Priority":     zog.Ptr(zog.Int32().GTE(1).LTE(5)),
	"StartedAt":    zog.Ptr(zog.Time()),
	"DurationDays": zog.Ptr(zog.Int32()),
	"FinishedAt":   zog.Ptr(zog.Time()),
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
	"StatusIDs": zog.Slice(zog.Int32()),
	"ManagerID": zog.Ptr(zog.Int32()),
	"StartAt":   zog.Ptr(zog.Time()),
	"EndAt":     zog.Ptr(zog.Time()),
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
	"ID":        zog.Int32().Required(),
	"NewName":   zog.String().Required(),
	"ManagerID": zog.Int32().Required(),
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
