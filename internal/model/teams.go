package model

type CreateTeamRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type UpdateTeamRequest struct {
	ID          int32  `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
}

type DeleteTeamRequest struct {
	ID int32 `json:"id" validate:"required"`
}
