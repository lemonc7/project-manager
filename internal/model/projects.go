package model

import (
	"fmt"
	"slices"
	"time"

	"github.com/Oudwins/zog"
	"github.com/Oudwins/zog/internals"
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
	"duration":    zog.Ptr(zog.Int32()),
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
	"name":        zog.Ptr(zog.String()),
	"description": zog.Ptr(zog.String()),
	"remark":      zog.Ptr(zog.String()),
	"status": zog.Ptr(zog.String().TestFunc(func(val *string, ctx internals.Ctx) bool {
		if val == nil {
			return true //字段可选，允许为nil
		}

		exists := slices.Contains([]string{"排单", "进行", "完成", "暂停", "终止"}, *val)
		if !exists {
			ctx.AddIssue(ctx.Issue().SetMessage(fmt.Sprintf("status(%s) is invalid", *val)))
			return false
		}

		return true
	})),
	"managerID": zog.Ptr(zog.Int32()),
	"priority":  zog.Ptr(zog.Int32().GTE(1).LTE(5)),
	"duration":  zog.Ptr(zog.Int32()),
})

type ListProjectsRequest struct {
	Status    []string   `json:"status"`
	ManagerID *int32     `json:"managerID"`
	StartAt   *time.Time `json:"startAt"`
	EndAt     *time.Time `json:"endAt"`
}

var ListProjectRequestSchema = zog.Struct(zog.Shape{
	"status":    zog.Slice(zog.String()).Default([]string{}),
	"managerID": zog.Ptr(zog.Int32()),
	"startAt":   zog.Ptr(zog.Time()),
	"endAt":     zog.Ptr(zog.Time()),
})

type CopyProjectRequest struct {
	NewName   string `json:"newName"`
	ManagerID int32  `json:"managerID"`
}

var CopyProjectRequestSchema = zog.Struct(zog.Shape{
	"newName":  zog.String().Required(),
	"mangerID": zog.Int32().Required(),
})

type IsProjectExistsRequest struct {
	NewName string `query:"newName"`
}

var IsProjectExistsRequestSchema = zog.Struct(zog.Shape{
	"newName": zog.String().Required(),
})
