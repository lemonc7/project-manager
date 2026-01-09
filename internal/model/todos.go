package model

import "github.com/Oudwins/zog"

type ListTodosRequest struct {
	Status []string `json:"status"`
}

var ListTodoRequestSchema = zog.Struct(zog.Shape{
	"status": zog.Slice(zog.String()).Default([]string{}),
})

type CreateTodoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Remark      string `json:"remark"`
	Priority    int32  `json:"priority"`
	ReporterID  int32  `json:"reporterID"`
	AssigneeID  int32  `json:"assigneeID"`
}

var CreateTodoRequestSchema = zog.Struct(zog.Shape{
	"name":        zog.String().Required(),
	"description": zog.String(),
	"remark":      zog.String(),
	"priority":    zog.Int32().Required().GTE(1).LTE(5),
	"reporterID":  zog.Int32().Required(),
	"assigneeID":  zog.Int32().Required(),
})

type UpdateTodoRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Remark      *string `json:"remark"`
	Priority    *int32  `json:"priority"`
	Status      *string `json:"status"`
	ReporterID  *int32  `json:"reporterID"`
	AssigneeID  *int32  `json:"assigneeID"`
}

var UpdateTodoRequestSchema = zog.Struct(zog.Shape{
	"name":        zog.Ptr(zog.String()),
	"description": zog.Ptr(zog.String()),
	"remark":      zog.Ptr(zog.String()),
	"priority":    zog.Ptr(zog.Int32().GTE(1).LTE(5)),
	"status":      zog.Ptr(zog.String()),
	"reporterID":  zog.Ptr(zog.Int32()),
	"assigneeID":  zog.Ptr(zog.Int32()),
})
