package model

import (
	"github.com/Oudwins/zog"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var loginRequestSchema = zog.Struct(zog.Shape{
	"username": zog.String().Required(),
	"password": zog.String().Required(),
})

func (r *LoginRequest) Validate() error {
	return SchemaValidate(loginRequestSchema, r)
}

type LoginResponse struct {
	AccessToken string   `json:"accessToken"`
	UserID      int32    `json:"userID"`
	Username    string   `json:"username"`
	Roles       []string `json:"roles,omitempty"`
}

type InitAdminRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var initAdminRequestSchema = zog.Struct(zog.Shape{
	"name":     zog.String().Required(),
	"username": zog.String().Required().Min(2).Max(20),
	"password": zog.String().Required().Min(6).Max(20),
})

func (r *InitAdminRequest) Validate() error {
	return SchemaValidate(initAdminRequestSchema, r)
}

type ResetPasswordRequest struct {
	ID       int32  `param:"id"`
	Password string `json:"password"`
}

var resetPasswordRequestSchema = zog.Struct(zog.Shape{
	"id":       zog.Int32().Required(),
	"password": zog.String().Required().Min(6).Max(20),
})

func (r *ResetPasswordRequest) Validate() error {
	return SchemaValidate(resetPasswordRequestSchema, r)
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

var changePasswordRequestSchema = zog.Struct(zog.Shape{
	"oldPassword": zog.String().Required(),
	"newPassword": zog.String().Required().Min(6).Max(20),
})

func (r *ChangePasswordRequest) Validate() error {
	return SchemaValidate(changePasswordRequestSchema, r)
}
