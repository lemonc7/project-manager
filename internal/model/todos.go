package model

import (
	"github.com/Oudwins/zog"
	"github.com/lemonc7/project-manager/pkg/util"
)

type ListTodosRequest struct {
	ID     int32    `param:"id"`
	Status []string `json:"status"`
}

var listTodosRequestSchema = zog.Struct(zog.Shape{
	"id":     zog.Int32().Required(),
	"status": zog.Slice(zog.String()).Default([]string{}),
})

func (r *ListTodosRequest) Validate() error {
	return util.SchemaValidate(listTodosRequestSchema, r)
}

type CreateTodoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Remark      string `json:"remark"`
	Priority    int32  `json:"priority"`
	ReporterID  int32  `json:"reporterID"`
	AssigneeID  int32  `json:"assigneeID"`
}

var createTodoRequestSchema = zog.Struct(zog.Shape{
	"name":        zog.String().Required(),
	"description": zog.String(),
	"remark":      zog.String(),
	"priority":    zog.Int32().Required().GTE(1).LTE(5),
	"reporterID":  zog.Int32().Required(),
	"assigneeID":  zog.Int32().Required(),
})

func (r *CreateTodoRequest) Validate() error {
	return util.SchemaValidate(createTodoRequestSchema, r)
}

type UpdateTodoRequest struct {
	ID          int32   `param:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Remark      *string `json:"remark"`
	Priority    *int32  `json:"priority"`
	Status      *string `json:"status"`
	ReporterID  *int32  `json:"reporterID"`
	AssigneeID  *int32  `json:"assigneeID"`
}

var updateTodoRequestSchema = zog.Struct(zog.Shape{
	"id":          zog.Int32().Required(),
	"name":        zog.Ptr(zog.String()),
	"description": zog.Ptr(zog.String()),
	"remark":      zog.Ptr(zog.String()),
	"priority":    zog.Ptr(zog.Int32().GTE(1).LTE(5)),
	"status":      zog.Ptr(zog.String()),
	"reporterID":  zog.Ptr(zog.Int32()),
	"assigneeID":  zog.Ptr(zog.Int32()),
})

func (r *UpdateTodoRequest) Validate() error {
	return util.SchemaValidate(updateTodoRequestSchema, r)
}
