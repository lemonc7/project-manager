package model

import (
	"github.com/Oudwins/zog"
	"github.com/lemonc7/project-manager/internal/repo"
	"github.com/lemonc7/project-manager/pkg/util"
)

type InsertProcessRequest struct {
	Sequence    int32  `param:"sequence"`
	Name        string `json:"name"`
	ProjectID   int32  `json:"projectID"`
	EmployeeID  int32  `json:"employeeID"`
	Description string `json:"description"`
	Remark      string `json:"remark"`
}

var insertProcessRequestSchema = zog.Struct(zog.Shape{
	"sequence":    zog.Int32().Required().GTE(1),
	"name":        zog.String().Required(),
	"projectID":   zog.Int32().Required(),
	"employeeID":  zog.Int32(),
	"description": zog.String(),
	"remark":      zog.String(),
})

func (r *InsertProcessRequest) Validate() error {
	return util.SchemaValidate(insertProcessRequestSchema, r)
}

type UpdateProcessRequest struct {
	ID          int32               `param:"id"`
	Name        *string             `json:"name"`
	Description *string             `json:"description"`
	Remark      *string             `json:"remark"`
	Status      *repo.ProcessStatus `json:"status"`
	EmployeeID  *int32              `json:"employeeID"`
	SpendTime   *float64            `json:"spendTime"`
}

var updateProcessRequestSchema = zog.Struct(zog.Shape{
	"id":          zog.Int32().Required(),
	"name":        zog.Ptr(zog.String()),
	"description": zog.Ptr(zog.String()),
	"remark":      zog.Ptr(zog.String()),
	"status": zog.Ptr(zog.String().OneOf([]string{
		"等待", "准备", "进行", "返工", "完成", "暂停", "取消"})),
	"employeeID": zog.Ptr(zog.Int32()),
	"spendTime":  zog.Ptr(zog.Float64()),
})

func (r *UpdateProcessRequest) Validate() error {
	return util.SchemaValidate(updateProcessRequestSchema, r)
}

type FinishProcessRequest struct {
	ID        int32   `param:"id"`
	ProjectID int32   `json:"projectID"`
	Sequence  int32   `json:"sequence"`
	SpendTime float64 `json:"spendTime"`
	Remark    *string `json:"remark"`
}

var finishProcessRequestSchema = zog.Struct(zog.Shape{
	"id":        zog.Int32().Required(),
	"projectID": zog.Int32().Required(),
	"sequence":  zog.Int32().Required(),
	"spendTime": zog.Float64().GTE(0),
	"remark":    zog.Ptr(zog.String()),
})

func (r *FinishProcessRequest) Validate() error {
	return util.SchemaValidate(finishProcessRequestSchema, r)
}

type ReworkProcessRequest struct {
	ID        int32   `param:"id"`
	SpendTime float64 `json:"spendTime"`
	Remark    string  `json:"remark"`
}

var reworkProcessRequestSchema = zog.Struct(zog.Shape{
	"id":        zog.Int32().Required(),
	"spendTime": zog.Float64().GTE(0),
	"remark":    zog.String(),
})

func (r *ReworkProcessRequest) Validate() error {
	return util.SchemaValidate(reworkProcessRequestSchema, r)
}

type DeleteProcessesRequest struct {
	ProjectID int32   `param:"projectID"`
	Sequences []int32 `json:"sequences"`
}

var deleteProcessesRequestSchema = zog.Struct(zog.Shape{
	"projectID": zog.Int32().Required(),
	"sequence":  zog.Slice(zog.Int32().Required().GTE(1)).Required(),
})

func (r *DeleteProcessesRequest) Validate() error {
	return util.SchemaValidate(deleteProcessesRequestSchema, r)
}

type ListTasksRequest struct {
	EmployeeID int32  `param:"employeeID"`
	Status     string `query:"status"`
}

var listTasksRequestSchema = zog.Struct(zog.Shape{
	"employeeID": zog.Int32().Required(),
	"status":     zog.String().Required().OneOf([]string{"todo", "ready", "doing"}),
})

func (r *ListTasksRequest) Validate() error {
	return util.SchemaValidate(listTasksRequestSchema, r)
}
