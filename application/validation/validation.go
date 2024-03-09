package validation

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

func (validator *Validator) ValidateStruct(structToBeValidated interface{}) error {
	return validator.validator.Struct(structToBeValidated)
}
