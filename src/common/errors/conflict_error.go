package errors

import "fmt"

// ConflictError maps to HTTP 409.
type ConflictError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewConflictError(code, message string) *ConflictError {
	return &ConflictError{Code: code, Message: message}
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
