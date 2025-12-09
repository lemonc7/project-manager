package model

type CreateEmployeeRequest struct {
	Name     string   `json:"name" validate:"required"`
	Username string   `json:"username" validate:"required,max=20"`
	Password string   `json:"password" validate:"required,min=6,max=20"`
	Roles    []string `json:"roles"`
	TeamID   *int32   `json:"teamID"`
}

type UpdateEmployeeRequest struct {
	Name     *string  `json:"name"`
	Roles    []string `json:"roles"`
	TeamID   *int32   `json:"teamID"`
	IsActive *bool    `json:"isActive"`
}
