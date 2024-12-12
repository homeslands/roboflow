package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

const (
	alphaNumberSpaceString = "^[a-zA-Z0-9 ]+$"
	qrCodeString           = "^[a-zA-Z0-9_-]+$"
)

var (
	AlphaNumberSpaceRegex = regexp.MustCompile(alphaNumberSpaceString)
	QRCodeRegex           = regexp.MustCompile(qrCodeString)
)

// AlphaNumberSpaceValidator checks if the given string contains only alphabets, numbers and spaces
func AlphaNumberSpaceValidator(fl validator.FieldLevel) bool {
	return AlphaNumberSpaceRegex.MatchString(fl.Field().String())
}

// QRCodeValidator checks if the given string contains only alphabets, numbers, underscores and hyphens
func QRCodeValidator(fl validator.FieldLevel) bool {
	return QRCodeRegex.MatchString(fl.Field().String())
}
