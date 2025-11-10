package model

import "github.com/lemonc7/project-manager/internal/repo"

type CreateTeamRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty" validate:"omitempty"`
}

type UpdateTeamRequest struct {
	ID          int32  `json:"name" validate:"required"`
	NewName     string `json:"new_name" validate:"required"`
	Description string `json:"description,omitempty" validate:"omitempty"`
}

type DeleteTeamRequest struct {
	ID int32 `json:"id" validate:"required"`
}

type ListTeamsResponse struct {
	Teams []repo.Team `json:"teams"`
}
