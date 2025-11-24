package model

type LoginRequest struct {
	Username string `json:"username" validate:"required,max=20"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken string   `json:"accessToken"`
	UserID      int32    `json:"userID"`
	Username    string   `json:"username"`
	Roles       []string `json:"roles,omitempty"`
}

type InitAdminRequest struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required,max=20"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type ResetPasswordRequest struct {
	ID       int32  `json:"id" validate:"required"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=6,max=20"`
}
