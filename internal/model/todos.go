package model

import (
	"time"

	"github.com/Oudwins/zog"
)

type ListTodosRequest struct {
	ID int32 `param:"id"`
}

func (r *ListTodosRequest) Validate() error {
	return SchemaValidate(idSchema, r)
}

type CreateTodoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Remark      string `json:"remark"`
	Priority    int32  `json:"priority"`
	CreatorID   int32  `json:"creatorID"`
	AssigneeID  int32  `json:"assigneeID"`
}

var createTodoRequestSchema = zog.Struct(zog.Shape{
	"name":        zog.String().Required(),
	"description": zog.String(),
	"remark":      zog.String(),
	"priority":    zog.Int32().Required().GTE(1).LTE(5),
	"creatorID":   zog.Int32().Required(),
	"assigneeID":  zog.Int32().Required(),
})

func (r *CreateTodoRequest) Validate() error {
	return SchemaValidate(createTodoRequestSchema, r)
}

type UpdateTodoRequest struct {
	ID          int32      `param:"id"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Remark      *string    `json:"remark"`
	Priority    *int32     `json:"priority"`
	StatusID    *int32     `json:"statusID"`
	CreatorID   *int32     `json:"creatorID"`
	AssigneeID  *int32     `json:"assigneeID"`
	StartedAt   *time.Time `json:"startedAt"`
	FinishedAt  *time.Time `json:"finishedAt"`
}

var updateTodoRequestSchema = zog.Struct(zog.Shape{
	"id":          zog.Int32().Required(),
	"name":        zog.Ptr(zog.String()),
	"description": zog.Ptr(zog.String()),
	"remark":      zog.Ptr(zog.String()),
	"priority":    zog.Ptr(zog.Int32().GTE(1).LTE(5)),
	"statusID":    zog.Ptr(zog.Int32()),
	"creatorID":   zog.Ptr(zog.Int32()),
	"assigneeID":  zog.Ptr(zog.Int32()),
	"startedAt":   zog.Ptr(zog.Time()),
	"finishedAt":  zog.Ptr(zog.Time()),
})

func (r *UpdateTodoRequest) Validate() error {
	return SchemaValidate(updateTodoRequestSchema, r)
}

type DeleteTodoRequest struct {
	ID int32 `param:"id"`
}

func (r *DeleteTodoRequest) Validate() error {
	return SchemaValidate(idSchema, r)
}
