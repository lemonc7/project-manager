package model

import (
	"github.com/Oudwins/zog"
)

type CreateTeamRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

var createTeamRequestSchema = zog.Struct(zog.Shape{
	"Name":        zog.String().Required(),
	"Description": zog.String(),
})

func (r *CreateTeamRequest) Validate() error {
	return SchemaValidate(createTeamRequestSchema, r)
}

type UpdateTeamRequest struct {
	ID          int32   `param:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

var updateTeamRequestSchema = zog.Struct(zog.Shape{
	"ID":          zog.Int32().Required(),
	"Name":        zog.Ptr(zog.String()),
	"Description": zog.Ptr(zog.String()),
})

func (r *UpdateTeamRequest) Validate() error {
	return SchemaValidate(updateTeamRequestSchema, r)
}

type DeleteTeamRequest struct {
	ID int32 `param:"id"`
}

func (r *DeleteTeamRequest) Validate() error {
	return SchemaValidate(idSchema, r)
}
