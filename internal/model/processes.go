package model

import (
	"github.com/Oudwins/zog"
	"github.com/lemonc7/project-manager/internal/repo"
)

type InsertProcessRequest struct {
	Name        string `json:"name"`
	ProjectID   int32  `json:"projectID"`
	EmployeeID  int32  `json:"employeeID"`
	Description string `json:"description"`
	Remark      string `json:"remark"`
}

var InsertProcessRequestSchema = zog.Struct(zog.Shape{
	"name":        zog.String().Required(),
	"projectID":   zog.Int32().Required(),
	"employeeID":  zog.Int32(),
	"description": zog.String(),
	"remark":      zog.String(),
})

type UpdateProcessRequest struct {
	Name        *string             `json:"name"`
	Description *string             `json:"description"`
	Remark      *string             `json:"remark"`
	Status      *repo.ProcessStatus `json:"status"`
	EmployeeID  *int32              `json:"employeeID"`
	SpendTime   *float64            `json:"spendTime"`
}

var UpdateProcessRequestSchema = zog.Struct(zog.Shape{
	"name":        zog.String(),
	"description": zog.String(),
	"remark":      zog.String(),
	"status":      zog.String(),
	"employeeID":  zog.Int32(),
	"spendTime":   zog.Float64(),
})

type FinishProcessRequest struct {
	ProjectID int32   `json:"projectID"`
	Sequence  int32   `json:"sequence"`
	SpendTime float64 `json:"spendTime"`
	Remark    *string `json:"remark"`
}

var FinishProcessRequestSchema = zog.Struct(zog.Shape{
	"projectID": zog.Int32().Required(),
	"sequence":  zog.Int32().Required(),
	"spendTime": zog.Float64().GTE(0),
	"remark":    zog.String(),
})

type ReworkProcessRequest struct {
	SpendTime float64 `json:"spendTime"`
	Remark    string  `json:"remark"`
}

var ReworkProcessRequestSchema = zog.Struct(zog.Shape{
	"spendTime": zog.Float64().GTE(0),
	"remark":    zog.String(),
})

type DeleteProcessesRequest struct {
	Sequences []int32 `json:"sequences"`
}

var DeleteProcessesRequestSchema = zog.Struct(zog.Shape{
	"sequences": zog.Slice(zog.Int32()).Required(),
})
