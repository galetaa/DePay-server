package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"transaction-validation-service/internal/controllers"
	"transaction-validation-service/internal/models"
	"transaction-validation-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupValidationRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	svc := services.NewValidationService()
	ctrl := controllers.NewValidationController(svc)

	router := gin.New()
	router.POST("/transaction/validate", ctrl.ValidateTransaction)
	return router
}

func TestValidateTransactionSuccess(t *testing.T) {
	router := setupValidationRouter()
	reqBody := models.ValidationRequest{
		TransactionID:    "tx123",
		SignedData:       "signature",
		SenderAddress:    "0x1234567890abcdef1234567890abcdef12345678", // 42 символа
		RecipientAddress: "0xabcdef1234567890abcdef1234567890abcdef12", // 42 символа
		Amount:           "1000000000000000000",
		Currency:         "ETH",
	}
	jsonValue, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/transaction/validate", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp["status"])
}

func TestValidateTransactionInsufficientFunds(t *testing.T) {
	router := setupValidationRouter()
	reqBody := models.ValidationRequest{
		TransactionID:    "tx123",
		SignedData:       "signature",
		SenderAddress:    "0x1234567890abcdef1234567890abcdef12345678",
		RecipientAddress: "0xabcdef1234567890abcdef1234567890abcdef12",
		Amount:           "0", // недостаточно средств
		Currency:         "ETH",
	}
	jsonValue, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/transaction/validate", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "failed", resp["status"])
	assert.Equal(t, "insufficient funds", resp["error"])
}

func TestValidateTransactionInvalidAddress(t *testing.T) {
	router := setupValidationRouter()
	reqBody := models.ValidationRequest{
		TransactionID:    "tx123",
		SignedData:       "signature",
		SenderAddress:    "invalid", // неверный формат
		RecipientAddress: "invalid",
		Amount:           "1000000000000000000",
		Currency:         "ETH",
	}
	jsonValue, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/transaction/validate", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "failed", resp["status"])
	assert.Equal(t, "invalid address format", resp["error"])
}
