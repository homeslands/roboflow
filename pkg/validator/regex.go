package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

const (
	alphaNumberSpaceString = "^[a-zA-Z0-9 ]+$"
)

var (
	alphaNumberSpaceRegex = regexp.MustCompile(alphaNumberSpaceString)
)

// alphaNumberSpaceValidator checks if the given string contains only alphabets, numbers and spaces
func alphaNumberSpaceValidator(fl validator.FieldLevel) bool {
	return alphaNumberSpaceRegex.MatchString(fl.Field().String())
}
