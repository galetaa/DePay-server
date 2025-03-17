package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"kyc-service/internal/controllers"
	"kyc-service/internal/models"
	"kyc-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	kycSvc := services.NewKYCService()
	kycCtrl := controllers.NewKYCController(kycSvc)
	router := gin.New()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.POST("/kyc", kycCtrl.ProcessKYC)
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

func TestKYCProcess(t *testing.T) {
	router := setupRouter()
	reqBody := models.KYCRequest{
		UserID:       "user123",
		DocumentType: "passport",
		DocumentData: "dummy-data",
	}
	jsonBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/kyc", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.KYCResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "user123", resp.UserID)
	assert.Equal(t, "verified", resp.KYCStatus)
	assert.Equal(t, "KYC verification successful", resp.Message)
}
