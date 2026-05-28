package models

import "fmt"

// APIError represents a structured API error
type APIError struct {
	Code    int
	Message string
	Err     error
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// BadRequest creates a 400 error
func BadRequest(msg string) *APIError {
	return &APIError{Code: 400, Message: msg}
}

// Unauthorized creates a 401 error
func Unauthorized(msg string) *APIError {
	return &APIError{Code: 401, Message: msg}
}

// NotFound creates a 404 error
func NotFound(msg string) *APIError {
	return &APIError{Code: 404, Message: msg}
}

// InternalError creates a 500 error
func InternalError(msg string, err error) *APIError {
	return &APIError{Code: 500, Message: msg, Err: err}
}
