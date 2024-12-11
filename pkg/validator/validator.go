package validator

import "github.com/go-playground/validator/v10"

// Validator is a wrapper around the go-playground/validator
type Validator struct {
	validate *validator.Validate
}

// New creates a new validator
func New() *Validator {
	v := validator.New(validator.WithRequiredStructEnabled())

	// Register custom validators here
	v.RegisterValidation("alphanumspace", alphaNumberSpaceValidator)
	return &Validator{
		validate: v,
	}
}

// Validate validates the given struct
func (v *Validator) Validate(s interface{}) error {
	return v.validate.Struct(s)
}

// IsValidationError checks if the given error is a validation error
func IsValidationError(err error) bool {
	_, ok := err.(validator.ValidationErrors)
	return ok
}
