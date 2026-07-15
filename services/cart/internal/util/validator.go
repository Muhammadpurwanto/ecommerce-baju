package util

import (
	"github.com/go-playground/validator/v10"
)

func FormatValidationErrors(err error) []map[string]string {
	var errors []map[string]string
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, map[string]string{
			"field":   e.Field(),
			"message": e.Tag() + " validation failed",
		})
	}
	return errors
}
