package utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) map[string]string {
	errors := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		errors[e.Field()] = fmt.Sprintf("не проходит '%s'", e.Tag())
	}
	return errors
}
