package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

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

// requireUserID returns the authenticated user id from the request context.
// When the value is missing or zero, it writes a 401 JSON response and returns
// the resulting echo error so the handler can return early.
func requireUserID(c *echo.Context) (int64, error) {
	if userID, ok := c.Get("user_id").(int64); ok && userID != 0 {
		return userID, nil
	}
	return 0, c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
}
