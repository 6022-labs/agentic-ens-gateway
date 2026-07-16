package errors

import "fmt"

// InternalError maps to HTTP 500.
type InternalError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewInternalError(code, message string) *InternalError {
	return &InternalError{Code: code, Message: message}
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
