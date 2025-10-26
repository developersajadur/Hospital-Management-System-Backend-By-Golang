package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    int
	Message string
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
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"success":    true,
		"statusCode": statusCode,
		"message":    message,
		"data":       data,
	})
}

// Error sends a standardized error response
func Error(c *gin.Context, err error) {
	// Default to 500
	status := http.StatusInternalServerError
	msg := err.Error()

	// If the error is an AppError, override status
	if appErr, ok := err.(*AppError); ok {
		status = appErr.Code
		msg = appErr.Message
	}

	c.JSON(status, gin.H{
		"success":    false,
		"statusCode": status,
		"message":    msg,
	})
}
