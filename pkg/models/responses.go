package models

// SuccessResponse wraps successful responses
type SuccessResponse struct {
	Data any `json:"data"`
}

// ErrorResponse wraps error responses
type ErrorResponse struct {
	Error string `json:"error"`
}

// MessageResponse for simple string responses
type MessageResponse struct {
	Message string `json:"message"`
}
