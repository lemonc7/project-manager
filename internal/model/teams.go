package model

type CreateTeamRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty" validate:"omitempty,max=255"`
}

type UpdateTeamRequest struct {
	
}