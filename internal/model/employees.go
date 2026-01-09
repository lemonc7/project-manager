package model

import "github.com/Oudwins/zog"

type CreateEmployeeRequest struct {
	Name     string   `json:"name"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
	TeamID   *int32   `json:"teamID"`
}

var CreateEmployeeSchema = zog.Struct(zog.Shape{
	"name":     zog.String().Required(),
	"username": zog.String().Required().Min(2).Max(20),
	"password": zog.String().Required().Min(6).Max(20),
	"roles":    zog.Slice(zog.String()).Default([]string{}),
	"teamID":   zog.Ptr(zog.Int32()),
})

type UpdateEmployeeRequest struct {
	Name     *string  `json:"name"`
	Roles    []string `json:"roles"`
	TeamID   *int32   `json:"teamID"`
	IsActive *bool    `json:"isActive"`
}

var UpdateEmployeeSchema = zog.Struct(zog.Shape{
	"name":     zog.Ptr(zog.String()),
	"roles":    zog.Slice(zog.String()),
	"teamID":   zog.Ptr(zog.Int32()),
	"isActive": zog.Ptr(zog.Bool()),
})
