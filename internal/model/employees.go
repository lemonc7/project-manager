package model

type CreateEmployeeRequest struct {
	Name   string `json:"name" validate:"required"`
	TeamID *int32 `json:"team_id"`
}

type UpdateEmployeeRequest struct {
	Name     *string `json:"name"`
	TeamID   *int32  `json:"team_id"`
	IsActive *bool   `json:"is_active"`
}
