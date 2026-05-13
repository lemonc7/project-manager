package model

import (
	"github.com/Oudwins/zog"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var loginRequestSchema = zog.Struct(zog.Shape{
	"Username": zog.String().Required(),
	"Password": zog.String().Required(),
})

func (r *LoginRequest) Validate() error {
	return SchemaValidate(loginRequestSchema, r)
}

type LoginResponse struct {
	AccessToken string  `json:"accessToken"`
	UserID      int32   `json:"userID"`
	Username    string  `json:"username"`
	RoleIDs     []int32 `json:"roleIDs,omitempty"`
}

type InitAdminRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var initAdminRequestSchema = zog.Struct(zog.Shape{
	"Name":     zog.String().Required(),
	"Username": zog.String().Required().Min(2).Max(20),
	"Password": zog.String().Required().Min(6).Max(20),
})

func (r *InitAdminRequest) Validate() error {
	return SchemaValidate(initAdminRequestSchema, r)
}

type ResetPasswordRequest struct {
	ID       int32  `param:"id"`
	Password string `json:"password"`
}

var resetPasswordRequestSchema = zog.Struct(zog.Shape{
	"ID":       zog.Int32().Required(),
	"Password": zog.String().Required().Min(6).Max(20),
})

func (r *ResetPasswordRequest) Validate() error {
	return SchemaValidate(resetPasswordRequestSchema, r)
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

var changePasswordRequestSchema = zog.Struct(zog.Shape{
	"OldPassword": zog.String().Required(),
	"NewPassword": zog.String().Required().Min(6).Max(20),
})

func (r *ChangePasswordRequest) Validate() error {
	return SchemaValidate(changePasswordRequestSchema, r)
}
