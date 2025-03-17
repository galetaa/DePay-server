package errors

import "fmt"

// AppError определяет тип ошибки приложения с кодом и сообщением
type AppError struct {
	Code    int
	Message string
}

// Error возвращает строковое представление ошибки
func (e *AppError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

// NewAppError создает новый экземпляр AppError
func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}
