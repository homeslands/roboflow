package validator

import "github.com/go-playground/validator/v10"

var (
	validate *validator.Validate
)

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())

	// Register custom validators here
	validate.RegisterValidation("alphanumspace", AlphaNumberSpaceValidator)
	validate.RegisterValidation("qrcode", QRCodeValidator)
}

// Validate validates the given struct
func Validate(s interface{}) error {
	return validate.Struct(s)
}

// IsValidationError checks if the given error is a validation error
func IsValidationError(err error) bool {
	_, ok := err.(validator.ValidationErrors)
	return ok
}
