package util

import (
	"github.com/go-playground/validator/v10"
)

func FormatValidationErrors(err error) []map[string]string {
	var errors []map[string]string
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, map[string]string{
			"field":   e.Field(),
			"message": getValidationMessage(e),
		})
	}
	return errors
}

func getValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return e.Field() + " is required"
	case "min":
		return e.Field() + " must be at least " + e.Param() + " characters"
	case "max":
		return e.Field() + " must be at most " + e.Param() + " characters"
	case "email":
		return e.Field() + " must be a valid email"
	case "url":
		return e.Field() + " must be a valid URL"
	default:
		return e.Field() + " is invalid"
	}
}
