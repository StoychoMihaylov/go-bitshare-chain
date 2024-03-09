package commands

import (
	validation "bitshare-chain/application/validation"
	context "context"
)

type TestCommand struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
}

type TestCommandHandler struct {
	validator *validation.Validator
}

func NewTestCommandHandler(validator *validation.Validator) *TestCommandHandler {
	return &TestCommandHandler{
		validator: validator,
	}
}

func (handler *TestCommandHandler) Handle(context context.Context, command TestCommand) (bool, error) {

	if err := handler.validator.ValidateStruct(command); err != nil {
		// Validation failed, return an error response
		return false, err
	}

	// Return a success response
	return true, nil
}
