package validation

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func Struct(s any) error {
	return wrapValidatorError(validate.Struct(s))
}

type Error struct {
	Messages []string
}

func (v Error) Error() string {
	return strings.Join(v.Messages, "; ")
}

func wrapValidatorError(errorFromValidator error) error {
	var validateErrs validator.ValidationErrors
	if errors.As(errorFromValidator, &validateErrs) {
		result := Error{}
		for _, err := range validateErrs {
			result.Messages = append(result.Messages, err.Error())
		}

		return result
	}

	return errorFromValidator
}
