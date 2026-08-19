package errors

import "fmt"

// UnauthorizedError maps to HTTP 401.
type UnauthorizedError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewUnauthorizedError(code, message string) *UnauthorizedError {
	return &UnauthorizedError{Code: code, Message: message}
}

func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
