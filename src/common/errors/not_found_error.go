package errors

import "fmt"

// NotFoundError maps to HTTP 404.
type NotFoundError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewNotFoundError(code, message string) *NotFoundError {
	return &NotFoundError{Code: code, Message: message}
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
