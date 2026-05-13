package model

import "github.com/go-playground/validator/v10"

type ValidateAdapter struct {
	Validate *validator.Validate
}

func (a *ValidateAdapter) ValidateStruct(ptr any) error {
	return a.Validate.Struct(ptr)
}
