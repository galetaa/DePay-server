package controllers

import (
	"net/http"
	"strings"
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
		c.JSON(statusCodeForTransactionError(err), gin.H{"error": err.Error()})
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
	if transactionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "transaction_id is required"})
		return
	}
	if err := tc.service.UpdateStatus(transactionID, "submitted", ""); err != nil {
		c.JSON(statusCodeForTransactionError(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"transaction_id": transactionID, "status": "submitted"}})
}

func (tc *TransactionController) ValidateTransaction(c *gin.Context) {
	transactionID := c.Param("transaction_id")
	if err := tc.service.UpdateStatus(transactionID, "validated", ""); err != nil {
		c.JSON(statusCodeForTransactionError(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"transaction_id": transactionID, "status": "validated"}})
}

func (tc *TransactionController) BroadcastTransaction(c *gin.Context) {
	transactionID := c.Param("transaction_id")
	tx, err := tc.service.Broadcast(transactionID)
	if err != nil {
		status := statusCodeForTransactionError(err)
		if status == http.StatusInternalServerError {
			status = http.StatusBadGateway
		}
		c.JSON(status, gin.H{"error": err.Error(), "data": tx})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tx})
}

func (tc *TransactionController) ConfirmTransaction(c *gin.Context) {
	transactionID := c.Param("transaction_id")
	if err := tc.service.UpdateStatus(transactionID, "confirmed", ""); err != nil {
		c.JSON(statusCodeForTransactionError(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"transaction_id": transactionID, "status": "confirmed"}})
}

func (tc *TransactionController) GetTransaction(c *gin.Context) {
	tx, err := tc.service.Get(c.Param("transaction_id"))
	if err != nil {
		c.JSON(statusCodeForTransactionError(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tx})
}

func (tc *TransactionController) GetStatus(c *gin.Context) {
	tx, err := tc.service.Get(c.Param("transaction_id"))
	if err != nil {
		c.JSON(statusCodeForTransactionError(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": models.TransactionStatusResponse{
		TransactionID:    tx.TransactionID,
		Status:           tx.Status,
		FailureReason:    tx.FailureReason,
		BlockchainTxHash: tx.BlockchainTxHash,
	}})
}

func (tc *TransactionController) CancelTransaction(c *gin.Context) {
	transactionID := c.Param("transaction_id")
	if err := tc.service.UpdateStatus(transactionID, "cancelled", "user cancelled"); err != nil {
		c.JSON(statusCodeForTransactionError(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"transaction_id": transactionID, "status": "cancelled"}})
}

func statusCodeForTransactionError(err error) int {
	if err == nil {
		return http.StatusOK
	}
	message := strings.ToLower(err.Error())
	if strings.Contains(message, "not found") || strings.Contains(message, "no rows") {
		return http.StatusNotFound
	}
	if strings.Contains(message, "invalid") || strings.Contains(message, "terminal") || strings.Contains(message, "must be validated") {
		return http.StatusConflict
	}
	if strings.Contains(message, "required") {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
