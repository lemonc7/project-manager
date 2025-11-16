package model

type CreateEmployeeRequest struct {
	Name   string `json:"name" validate:"required"`
	TeamID *int32 `json:"teamID"`
}

type UpdateEmployeeRequest struct {
	Name     *string `json:"name"`
	TeamID   *int32  `json:"teamID"`
	IsActive *bool   `json:"isActive"`
}
