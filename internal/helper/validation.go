package helper

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) map[string]string {

	errors := map[string]string{}

	validationErrors, ok := err.(validator.ValidationErrors)

	if !ok {
		errors["error"] = err.Error()
		return errors
	}

	for _, fieldErr := range validationErrors {

		field := strings.ToLower(fieldErr.Field())

		switch fieldErr.Tag() {

		case "required":
			errors[field] = "field is required"

		case "min":
			errors[field] = "minimum length is " + fieldErr.Param()

		case "max":
			errors[field] = "maximum length is " + fieldErr.Param()

		case "gt":
			errors[field] = "must be greater than " + fieldErr.Param()

		case "oneof":
			errors[field] = "must be one of: " + fieldErr.Param()

		default:
			errors[field] = fieldErr.Error()
		}
	}

	return errors
}
