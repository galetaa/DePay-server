package controllers

import (
	"net/http"

	"transaction-validation-service/internal/models"
	"transaction-validation-service/internal/services"

	"github.com/gin-gonic/gin"
)

type ValidationController struct {
	service services.ValidationService
}

func NewValidationController(s services.ValidationService) *ValidationController {
	return &ValidationController{service: s}
}

// ValidateTransaction принимает запрос на валидацию и возвращает результат проверки.
func (vc *ValidationController) ValidateTransaction(c *gin.Context) {
	var req models.ValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	err := vc.service.Validate(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
