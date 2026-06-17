package handler

// APIError describes a client-visible error payload.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse builds a structured error response body.
func ErrorResponse(code, message string) APIError {
	return APIError{
		Code:    code,
		Message: message,
	}
}
