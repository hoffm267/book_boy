package errors

import "fmt"

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func ErrInvalidInput(msg string) error {
	return ValidationError{Message: msg}
}

func IsValidationError(err error) bool {
	_, ok := err.(ValidationError)
	return ok
}

var (
	ErrNotFound     = fmt.Errorf("resource not found")
	ErrUnauthorized = fmt.Errorf("unauthorized")
	ErrConflict     = fmt.Errorf("resource conflict")
)
