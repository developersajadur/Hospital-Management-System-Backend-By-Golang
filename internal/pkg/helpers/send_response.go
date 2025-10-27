package helpers

import (
	"encoding/json"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

// Helper to create new AppError
func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Success sends a standardized success response
func Success(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"statusCode": statusCode,
		"message":    message,
		"data":       data,
	})
}

// Error sends a standardized error response
func Error(w http.ResponseWriter, err error) {
	// Default to 500
	status := http.StatusInternalServerError
	msg := err.Error()

	// If the error is an AppError, override status
	if appErr, ok := err.(*AppError); ok {
		status = appErr.Code
		msg = appErr.Message
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    false,
		"statusCode": status,
		"message":    msg,
	})
}