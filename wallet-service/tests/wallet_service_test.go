package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"wallet-service/internal/controllers"
	"wallet-service/internal/models"
	"wallet-service/internal/repositories"
	"wallet-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	repo := repositories.NewWalletRepository()
	svc := services.NewWalletService(repo)
	ctrl := controllers.NewWalletController(svc)

	router := gin.New()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/wallet/export", ctrl.ExportWallets)
	router.POST("/wallet/balance", ctrl.GetBalance)
	return router
}

func TestHealthEndpoint(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp["status"])
}

func TestExportWallets(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/wallet/export", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var wallets []models.Wallet
	err := json.Unmarshal(w.Body.Bytes(), &wallets)
	assert.NoError(t, err)
	assert.Greater(t, len(wallets), 0, "Should return at least one wallet")
}

func TestGetBalance(t *testing.T) {
	router := setupRouter()
	reqBody := models.BalanceRequest{
		Address: "0x1234567890abcdef",
	}
	jsonBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/wallet/balance", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp models.BalanceResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "0x1234567890abcdef", resp.Address)
	assert.Equal(t, "ethereum", resp.Blockchain)
	assert.Equal(t, "1000000000000000000", resp.Balance)
}
