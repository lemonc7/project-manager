package model

type CreateEmployeeRequest struct {
	Name   string `json:"name" validate:"required"`
	TeamID *int32 `json:"team_id"`
}

type UpdateEmployeeRequest struct {
	ID       int32   `json:"id" validate:"required"`
	Name     *string `json:"name"`
	TeamID   *int32  `json:"team_id"`
	IsActive *bool   `json:"is_active"`
}

type DeleteEmployeeRequest struct {
	ID int32 `json:"id" validate:"required"`
}
