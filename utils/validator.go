package utils

import (
	"github.com/go-playground/validator/v10"
)

// CustomValidator wraps go-playground/validator for use with Echo.
type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}

// Validate runs struct validation. Echo calls this via c.Validate(...).
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// FormatValidationErrors turns validator errors into a readable field->message map.
func FormatValidationErrors(err error) map[string]string {
	errorsMap := make(map[string]string)

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		errorsMap["error"] = err.Error()
		return errorsMap
	}

	for _, fe := range validationErrors {
		field := fe.Field()
		switch fe.Tag() {
		case "required":
			errorsMap[field] = field + " is required"
		case "email":
			errorsMap[field] = "Must be a valid email address"
		case "min":
			errorsMap[field] = field + " must be at least " + fe.Param() + " characters"
		case "max":
			errorsMap[field] = field + " must be at most " + fe.Param() + " characters"
		case "gt":
			errorsMap[field] = field + " must be greater than " + fe.Param()
		case "oneof":
			errorsMap[field] = field + " must be one of: " + fe.Param()
		default:
			errorsMap[field] = "Invalid value for " + field
		}
	}
	return errorsMap
}
