package dto

// APIError represents a structured API error
type APIError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
}

// Error satisfies the error interface
func (e *APIError) Error() string {
	return e.Message
}

// NewAPIError creates a new AppError
func NewAPIError(statusCode int, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
	}
}
