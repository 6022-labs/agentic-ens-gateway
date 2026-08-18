package errors

import "fmt"

// UnavailableError maps to HTTP 502 (an upstream/dependency is unavailable).
type UnavailableError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewUnavailableError(code, message string) *UnavailableError {
	return &UnavailableError{Code: code, Message: message}
}

func (e *UnavailableError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
