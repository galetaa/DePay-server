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

func (wc *WalletController) CreateWallet(c *gin.Context) {
	var req models.CreateWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	wallet, err := wc.service.CreateWallet(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": wallet})
}

func (wc *WalletController) GetWallet(c *gin.Context) {
	wallet, err := wc.service.GetWallet(c.Param("wallet_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": wallet})
}

func (wc *WalletController) DeleteWallet(c *gin.Context) {
	if err := wc.service.DeleteWallet(c.Param("wallet_id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (wc *WalletController) GetWalletBalances(c *gin.Context) {
	balances, err := wc.service.GetWalletBalances(c.Param("wallet_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": balances})
}

func (wc *WalletController) SyncWallet(c *gin.Context) {
	c.JSON(http.StatusAccepted, gin.H{"data": gin.H{"status": "sync_scheduled", "wallet_id": c.Param("wallet_id")}})
}

func (wc *WalletController) GetBalanceByAddress(c *gin.Context) {
	address := c.Param("address")
	if address == "" {
		address = c.Param("wallet_id")
	}
	resp, err := wc.service.GetBalance(models.BalanceRequest{Address: address})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
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
