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
	"StepOrder":   zog.Int32().Required().GTE(1),
	"Name":        zog.String().Required(),
	"ProjectID":   zog.Int32().Required(),
	"EmployeeID":  zog.Int32().Required(),
	"Description": zog.String(),
	"Remark":      zog.String(),
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
	"ID":          zog.Int32().Required(),
	"Name":        zog.Ptr(zog.String()),
	"Description": zog.Ptr(zog.String()),
	"Remark":      zog.Ptr(zog.String()),
	"StatusID":    zog.Ptr(zog.Int32()),
	"EmployeeID":  zog.Ptr(zog.Int32()),
	"SpendHours":  zog.Ptr(zog.Float64()),
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
	"ID":         zog.Int32().Required(),
	"ProjectID":  zog.Int32().Required(),
	"StepOrder":  zog.Int32().Required(),
	"SpendHours": zog.Float64().GTE(0),
	"Remark":     zog.Ptr(zog.String()),
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
	"ID":         zog.Int32().Required(),
	"SpendHours": zog.Float64().GTE(0),
	"Remark":     zog.String(),
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
	"ProjectID":  zog.Int32().Required(),
	"StepOrders": zog.Slice(zog.Int32().Required().GTE(1)).Required(),
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
	EmployeeID int32               `param:"id"`
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
	"EmployeeID": zog.Int32().Required(),
	"Status":     SchemaStatusCategory.Required(),
})

func (r *ListTasksRequest) Validate() error {
	return SchemaValidate(listTasksRequestSchema, r)
}
