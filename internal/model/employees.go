package model

import (
	"github.com/Oudwins/zog"
)

type CreateEmployeeRequest struct {
	Name     string  `json:"name"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	RoleIDs  []int32 `json:"roleIDs"`
	TeamID   *int32  `json:"teamID"`
}

var createEmployeeSchema = zog.Struct(zog.Shape{
	"Name":     zog.String().Required(),
	"Username": zog.String().Required().Min(2).Max(20),
	"Password": zog.String().Required().Min(6).Max(20),
	"RoleIDs":  zog.Slice(zog.Int32()),
	"TeamID":   zog.Ptr(zog.Int32()),
})

func (r *CreateEmployeeRequest) Validate() error {
	return SchemaValidate(createEmployeeSchema, r)
}

type UpdateEmployeeRequest struct {
	ID       int32   `param:"id"`
	Name     *string `json:"name"`
	Username *string `json:"username"`
	RoleIDs  []int32 `json:"roleIDs"`
	TeamID   *int32  `json:"teamID"`
	IsActive *bool   `json:"isActive"`
}

var updateEmployeeSchema = zog.Struct(zog.Shape{
	"ID":       zog.Int32().Required(),
	"Name":     zog.Ptr(zog.String()),
	"Username": zog.Ptr(zog.String()),
	"RoleIDs":  zog.Slice(zog.Int32()),
	"TeamID":   zog.Ptr(zog.Int32()),
	"IsActive": zog.Ptr(zog.Bool()),
})

func (r *UpdateEmployeeRequest) Validate() error {
	return SchemaValidate(updateEmployeeSchema, r)
}

type DeleteEmployeeRequest struct {
	ID int32 `param:"id"`
}

var idSchema = zog.Struct(zog.Shape{
	"ID": zog.Int32().Required(),
})

func (r *DeleteEmployeeRequest) Validate() error {
	return SchemaValidate(idSchema, r)
}
