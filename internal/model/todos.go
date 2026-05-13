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
	"Name":        zog.String().Required(),
	"Description": zog.String(),
	"Remark":      zog.String(),
	"priority":    zog.Int32().Required().GTE(1).LTE(5),
	"CreatorID":   zog.Int32().Required(),
	"AssigneeID":  zog.Int32().Required(),
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
	"ID":          zog.Int32().Required(),
	"Name":        zog.Ptr(zog.String()),
	"Description": zog.Ptr(zog.String()),
	"Remark":      zog.Ptr(zog.String()),
	"Priority":    zog.Ptr(zog.Int32().GTE(1).LTE(5)),
	"StatusID":    zog.Ptr(zog.Int32()),
	"CreatorID":   zog.Ptr(zog.Int32()),
	"AssigneeID":  zog.Ptr(zog.Int32()),
	"StartedAt":   zog.Ptr(zog.Time()),
	"FinishedAt":  zog.Ptr(zog.Time()),
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
