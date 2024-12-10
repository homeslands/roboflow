package validator

import "github.com/go-playground/validator/v10"

type Validator interface {
	Validate(s interface{}) error
}

// xValidator is a wrapper around the go-playground/validator
type xValidator struct {
	validate *validator.Validate
}

// New creates a new validator
func New() *xValidator {
	return &xValidator{
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

// Validate validates the given struct
func (v *xValidator) Validate(s interface{}) error {
	return v.validate.Struct(s)
}

// IsValidationError checks if the given error is a validation error
func IsValidationError(err error) bool {
	_, ok := err.(validator.ValidationErrors)
	return ok
}
