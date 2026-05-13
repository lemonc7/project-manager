package model

import (
	"time"
)

type ListTodosRequest struct {
	ID int32 `param:"id" validate:"required"`
}

type CreateTodoRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Remark      string `json:"remark"`
	Priority    int32  `json:"priority" validate:"required,gte=1,lte=5"`
	CreatorID   int32  `json:"creatorID" validate:"required"`
	AssigneeID  int32  `json:"assigneeID" validate:"required"`
}

type UpdateTodoRequest struct {
	ID          int32      `param:"id" validate:"required"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Remark      *string    `json:"remark"`
	Priority    *int32     `json:"priority" validate:"gte=1,lte=5"`
	StatusID    *int32     `json:"statusID"`
	CreatorID   *int32     `json:"creatorID"`
	AssigneeID  *int32     `json:"assigneeID"`
	StartedAt   *time.Time `json:"startedAt"`
	FinishedAt  *time.Time `json:"finishedAt"`
}

type DeleteTodoRequest struct {
	ID int32 `param:"id" validate:"required"`
}
