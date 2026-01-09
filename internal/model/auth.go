package model

import "github.com/Oudwins/zog"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var LoginRequestSchema = zog.Struct(zog.Shape{
	"username": zog.String().Required(),
	"password": zog.String().Required(),
})

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

var InitAdminRequestSchema = zog.Struct(zog.Shape{
	"name":     zog.String().Required(),
	"username": zog.String().Required().Min(2).Max(20),
	"password": zog.String().Required().Min(6).Max(20),
})

type ResetPasswordRequest struct {
	Password string `json:"password"`
}

var ResetPasswordRequestSchema = zog.Struct(zog.Shape{
	"password": zog.String().Required().Min(6).Max(20),
})

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

var ChangePasswordRequestSchema = zog.Struct(zog.Shape{
	"oldPassword": zog.String().Required(),
	"newPassword": zog.String().Required().Min(6).Max(20),
})
