package domain

import "fmt"

var (
	ErrModelNotFound      = fmt.Errorf("model not found")
	ErrModelAlreadyExists = fmt.Errorf("model already exists")
	ErrInvalidModel       = fmt.Errorf("invalid model")
	ErrInvalidField       = fmt.Errorf("invalid field")
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}
