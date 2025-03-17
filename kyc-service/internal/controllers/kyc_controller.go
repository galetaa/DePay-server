package controllers

import (
	"net/http"

	"kyc-service/internal/models"
	"kyc-service/internal/services"

	"github.com/gin-gonic/gin"
)

type KYCController struct {
	service services.KYCService
}

func NewKYCController(s services.KYCService) *KYCController {
	return &KYCController{service: s}
}

// ProcessKYC обрабатывает KYC запрос, имитируя внешнюю интеграцию
func (kc *KYCController) ProcessKYC(c *gin.Context) {
	var req models.KYCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	resp, err := kc.service.ProcessKYC(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
