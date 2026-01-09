package model

import "github.com/Oudwins/zog"

type CreateTeamRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

var CreateTeamRequestSchema = zog.Struct(zog.Shape{
	"name":        zog.String().Required(),
	"description": zog.String(),
})

type UpdateTeamRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

var UpdateTeamRequestSchema = zog.Struct(zog.Shape{
	"name":        zog.Ptr(zog.String()),
	"description": zog.Ptr(zog.String()),
})
