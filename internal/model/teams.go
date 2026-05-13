package model

type CreateTeamRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type UpdateTeamRequest struct {
	ID          int32   `param:"id" validate:"required"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type DeleteTeamRequest struct {
	ID int32 `param:"id" validate:"required"`
}
