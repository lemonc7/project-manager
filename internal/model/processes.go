package model

import (
	"github.com/Oudwins/zog"
	"github.com/lemonc7/project-manager/internal/repo"
)

type InsertProcessRequest struct {
	StepOrder   int32  `param:"stepOrder"`
	Name        string `json:"name"`
	ProjectID   int32  `json:"projectID"`
	EmployeeID  int32  `json:"employeeID"`
	Description string `json:"description"`
	Remark      string `json:"remark"`
}

var insertProcessRequestSchema = zog.Struct(zog.Shape{
	"stepOrder":   zog.Int32().Required().GTE(1),
	"name":        zog.String().Required(),
	"projectID":   zog.Int32().Required(),
	"employeeID":  zog.Int32(),
	"description": zog.String(),
	"remark":      zog.String(),
})

func (r *InsertProcessRequest) Validate() error {
	return SchemaValidate(insertProcessRequestSchema, r)
}

type UpdateProcessRequest struct {
	ID          int32    `param:"id"`
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Remark      *string  `json:"remark"`
	StatusID    *int32   `json:"statusID"`
	EmployeeID  *int32   `json:"employeeID"`
	SpendHours  *float64 `json:"spendHours"`
}

var updateProcessRequestSchema = zog.Struct(zog.Shape{
	"id":          zog.Int32().Required(),
	"name":        zog.Ptr(zog.String()),
	"description": zog.Ptr(zog.String()),
	"remark":      zog.Ptr(zog.String()),
	"statusID":    zog.Ptr(zog.Int32()),
	"employeeID":  zog.Ptr(zog.Int32()),
	"spendHours":  zog.Ptr(zog.Float64()),
})

func (r *UpdateProcessRequest) Validate() error {
	return SchemaValidate(updateProcessRequestSchema, r)
}

type FinishProcessRequest struct {
	ID         int32   `param:"id"`
	ProjectID  int32   `json:"projectID"`
	StepOrder  int32   `json:"stepOrder"`
	SpendHours float64 `json:"spendHours"`
	Remark     *string `json:"remark"`
}

var finishProcessRequestSchema = zog.Struct(zog.Shape{
	"id":         zog.Int32().Required(),
	"projectID":  zog.Int32().Required(),
	"stepOrder":  zog.Int32().Required(),
	"spendHours": zog.Float64().GTE(0),
	"remark":     zog.Ptr(zog.String()),
})

func (r *FinishProcessRequest) Validate() error {
	return SchemaValidate(finishProcessRequestSchema, r)
}

type ReworkProcessRequest struct {
	ID         int32   `param:"id"`
	SpendHours float64 `json:"spendHours"`
	Remark     string  `json:"remark"`
}

var reworkProcessRequestSchema = zog.Struct(zog.Shape{
	"id":         zog.Int32().Required(),
	"spendHours": zog.Float64().GTE(0),
	"remark":     zog.String(),
})

func (r *ReworkProcessRequest) Validate() error {
	return SchemaValidate(reworkProcessRequestSchema, r)
}

type DeleteProcessRequest struct {
	ID int32 `param:"id"`
}

func (r *DeleteProcessRequest) Validate() error {
	return SchemaValidate(idSchema, r)
}

type DeleteProcessesRequest struct {
	ProjectID  int32   `param:"projectID"`
	StepOrders []int32 `json:"stepOrders"`
}

var deleteProcessesRequestSchema = zog.Struct(zog.Shape{
	"projectID":  zog.Int32().Required(),
	"stepOrders": zog.Slice(zog.Int32().Required().GTE(1)).Required(),
})

func (r *DeleteProcessesRequest) Validate() error {
	return SchemaValidate(deleteProcessesRequestSchema, r)
}

type ListProcessesRequest struct {
	ID int32 `param:"id"`
}

func (r *ListProcessesRequest) Validate() error {
	return SchemaValidate(idSchema, r)
}

type ListTasksRequest struct {
	EmployeeID int32               `param:"employeeID"`
	Status     repo.StatusCategory `query:"status"`
}

var SchemaStatusCategory = zog.StringLike[repo.StatusCategory]().
	OneOf([]repo.StatusCategory{
		repo.StatusCategoryTodo,
		repo.StatusCategoryReady,
		repo.StatusCategoryDoing,
		repo.StatusCategoryFinished,
	})

var listTasksRequestSchema = zog.Struct(zog.Shape{
	"employeeID": zog.Int32().Required(),
	"status":     SchemaStatusCategory.Required(),
})

func (r *ListTasksRequest) Validate() error {
	return SchemaValidate(listTasksRequestSchema, r)
}
