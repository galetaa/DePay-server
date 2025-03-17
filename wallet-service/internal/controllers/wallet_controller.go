package controllers

import (
	"net/http"
	"wallet-service/internal/models"
	"wallet-service/internal/services"

	"github.com/gin-gonic/gin"
)

type WalletController struct {
	service services.WalletService
}

func NewWalletController(s services.WalletService) *WalletController {
	return &WalletController{service: s}
}

// ExportWallets возвращает список кошельков
func (wc *WalletController) ExportWallets(c *gin.Context) {
	wallets, err := wc.service.ExportWallets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wallets)
}

// GetBalance принимает адрес кошелька и возвращает его баланс
func (wc *WalletController) GetBalance(c *gin.Context) {
	var req models.BalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	resp, err := wc.service.GetBalance(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
