package dto

type APIError struct {
	StatusCode int               `json:"-"`                // Exclude this field from the JSON response
	Details    map[string]string `json:"errors,omitempty"` // Hold validation errors only
	Message    string            `json:"-"`                // Exclude generic messages when returning errors
}

type APIErrorMap map[string]*APIError

// Error satisfies the error interface
func (e *APIError) Error() string {
	return e.Message
}

// NewAPIError creates a new AppError

func NewAPIError(statusCode int, message string, details map[string]string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		Details:    details,
	}
}
