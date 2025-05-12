package models

import (
	"fmt"
	"net/http"
)

// APIError represents a standardized API error
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Common error types
var (
	ErrNotFound       = APIError{Code: http.StatusNotFound, Message: "Resource not found"}
	ErrBadRequest     = APIError{Code: http.StatusBadRequest, Message: "Invalid request"}
	ErrInternalServer = APIError{Code: http.StatusInternalServerError, Message: "Internal server error"}
)

// Error implements the error interface
func (e APIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Details)
	}
	return e.Message
}