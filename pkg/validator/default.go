package validator

import (
	"github.com/go-playground/validator/v10"
)

// Validator is a validator that validates the given struct.
type Validator interface {
	// Validate validates the given struct
	Validate(s any) error
}

type defaultValidator struct {
	v10 *validator.Validate
}

// NewValidator creates a new default validator.
//
//nolint:revive
func NewValidator() *defaultValidator {
	return &defaultValidator{
		v10: newValidator(),
	}
}

func (v *defaultValidator) Validate(s any) error {
	return v.v10.Struct(s)
}
