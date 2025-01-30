package dto

type APIError struct {
	StatusCode int               `json:"-"`
	Details    map[string]string `json:"errors,omitempty"`
}

// NewAPIError creates a new AppError
func NewAPIError(statusCode int, details map[string]string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Details:    details,
	}
}
