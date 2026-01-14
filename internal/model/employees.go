package model

import (
	"github.com/Oudwins/zog"
	"github.com/lemonc7/project-manager/pkg/util"
)

type CreateEmployeeRequest struct {
	Name     string   `json:"name"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
	TeamID   *int32   `json:"teamID"`
}

var createEmployeeSchema = zog.Struct(zog.Shape{
	"name":     zog.String().Required(),
	"username": zog.String().Required().Min(2).Max(20),
	"password": zog.String().Required().Min(6).Max(20),
	"roles":    zog.Slice(zog.String()).Default([]string{}),
	"teamID":   zog.Ptr(zog.Int32()),
})

func (r *CreateEmployeeRequest) Validate() error {
	return util.SchemaValidate(createEmployeeSchema, r)
}

type UpdateEmployeeRequest struct {
	ID       int32    `param:"id"`
	Name     *string  `json:"name"`
	Roles    []string `json:"roles"`
	TeamID   *int32   `json:"teamID"`
	IsActive *bool    `json:"isActive"`
}

var updateEmployeeSchema = zog.Struct(zog.Shape{
	"id":       zog.Int32().Required(),
	"name":     zog.Ptr(zog.String()),
	"roles":    zog.Slice(zog.String()),
	"teamID":   zog.Ptr(zog.Int32()),
	"isActive": zog.Ptr(zog.Bool()),
})

func (r *UpdateEmployeeRequest) Validate() error {
	return util.SchemaValidate(updateEmployeeSchema, r)
}
