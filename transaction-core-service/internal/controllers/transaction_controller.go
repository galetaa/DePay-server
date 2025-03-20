package controllers

import (
	"net/http"
	"time"

	"transaction-core-service/internal/models"
	"transaction-core-service/internal/services"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	service services.TransactionService
}

func NewTransactionController(s services.TransactionService) *TransactionController {
	return &TransactionController{service: s}
}

// InitiateTransaction принимает транзакцию и инициирует её обработку
func (tc *TransactionController) InitiateTransaction(c *gin.Context) {
	var tx models.Transaction
	if err := c.ShouldBindJSON(&tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	// Если Timestamp не задан, задаем текущее время
	if tx.Timestamp.IsZero() {
		tx.Timestamp = time.Now().UTC()
	}
	// Валюта устанавливается в "ETH"
	tx.Currency = "ETH"

	if err := tc.service.Initiate(tx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "transaction initiated"})
}
