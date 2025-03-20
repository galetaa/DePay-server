package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gas-info-service/internal/controllers"
	"gas-info-service/internal/models"
	"gas-info-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupGasInfoRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	svc := services.NewGasService()
	ctrl := controllers.NewGasController(svc)

	router := gin.New()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/gas-info", ctrl.GetGasInfo)
	return router
}

func TestGasInfoHealth(t *testing.T) {
	router := setupGasInfoRouter()
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp["status"])
}

func TestGetGasInfo(t *testing.T) {
	router := setupGasInfoRouter()
	req, _ := http.NewRequest("GET", "/gas-info?network=polygon", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var info models.GasInfo
	err := json.Unmarshal(w.Body.Bytes(), &info)
	assert.NoError(t, err)
	assert.Equal(t, "polygon", info.Network)
	// Заглушка возвращает фиксированные значения: GasPrice: 50.0, EstimatedTime: 30, NetworkStatus: "normal"
	assert.Equal(t, 50.0, info.GasPrice)
	assert.Equal(t, 30, info.EstimatedTime)
	assert.Equal(t, "normal", info.NetworkStatus)
}
