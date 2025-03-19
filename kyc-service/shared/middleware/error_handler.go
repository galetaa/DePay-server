package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandlerMiddleware перехватывает ошибки, накопленные во время запроса
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Если были ошибки – возвращаем их в виде JSON
		if len(c.Errors) > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"errors": c.Errors})
		}
	}
}
