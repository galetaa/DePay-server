package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AppError определяет тип ошибки приложения с кодом и сообщением
type AppError struct {
	Code       int
	Message    string
	HTTPStatus int
	Details    map[string]any
}

// Error возвращает строковое представление ошибки
func (e *AppError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

// NewAppError создает новый экземпляр AppError
func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
	}
}

func NewAPIError(code int, message string, httpStatus int, details map[string]any) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Details:    details,
	}
}

func JSON(c *gin.Context, status int, code string, message string, details map[string]any) {
	if details == nil {
		details = map[string]any{}
	}
	c.JSON(status, gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
			"details": details,
		},
	})
}

func Data(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{"data": data})
}
