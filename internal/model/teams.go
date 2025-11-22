package model

type CreateTeamRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
}

type UpdateTeamRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
