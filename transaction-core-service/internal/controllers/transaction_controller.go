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
	if tx.Currency == "" {
		tx.Currency = "ETH"
	}
	if tx.Status == "" {
		tx.Status = "created"
	}

	if err := tc.service.Initiate(tx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "transaction initiated"})
}

func (tc *TransactionController) SubmitTransaction(c *gin.Context) {
	transactionID := c.Param("transaction_id")
	if transactionID == "" {
		var tx models.Transaction
		if err := c.ShouldBindJSON(&tx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		transactionID = tx.TransactionID
	}
	if err := tc.service.UpdateStatus(transactionID, "submitted", ""); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"transaction_id": transactionID, "status": "submitted"}})
}

func (tc *TransactionController) GetTransaction(c *gin.Context) {
	tx, err := tc.service.Get(c.Param("transaction_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tx})
}

func (tc *TransactionController) GetStatus(c *gin.Context) {
	tx, err := tc.service.Get(c.Param("transaction_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": models.TransactionStatusResponse{
		TransactionID: tx.TransactionID,
		Status:        tx.Status,
		FailureReason: tx.FailureReason,
	}})
}

func (tc *TransactionController) CancelTransaction(c *gin.Context) {
	transactionID := c.Param("transaction_id")
	if err := tc.service.UpdateStatus(transactionID, "cancelled", "user cancelled"); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"transaction_id": transactionID, "status": "cancelled"}})
}
