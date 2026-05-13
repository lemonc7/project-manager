package model

type CreateEmployeeRequest struct {
	Name     string  `json:"name" validate:"required"`
	Username string  `json:"username" validate:"required,min=2,max=20"`
	Password string  `json:"password" validate:"required,min=6,max=20"`
	RoleIDs  []int32 `json:"roleIDs" validate:"required"`
	TeamID   *int32  `json:"teamID"`
}

type UpdateEmployeeRequest struct {
	ID       int32   `param:"id" validate:"required"`
	Name     *string `json:"name"`
	Username *string `json:"username"`
	RoleIDs  []int32 `json:"roleIDs"`
	TeamID   *int32  `json:"teamID"`
	IsActive *bool   `json:"isActive"`
}

type DeleteEmployeeRequest struct {
	ID int32 `param:"id" validate:"required"`
}
