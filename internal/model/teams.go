package model

import (
	"github.com/Oudwins/zog"
	"github.com/lemonc7/project-manager/pkg/util"
)

type CreateTeamRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

var createTeamRequestSchema = zog.Struct(zog.Shape{
	"name":        zog.String().Required(),
	"description": zog.String(),
})

func (r *CreateTeamRequest) Validate() error {
	return util.SchemaValidate(createTeamRequestSchema, r)
}

type UpdateTeamRequest struct {
	ID          int32   `param:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

var updateTeamRequestSchema = zog.Struct(zog.Shape{
	"id":          zog.Int32().Required(),
	"name":        zog.Ptr(zog.String()),
	"description": zog.Ptr(zog.String()),
})

func (r *UpdateTeamRequest) Validate() error {
	return util.SchemaValidate(updateTeamRequestSchema, r)
}
