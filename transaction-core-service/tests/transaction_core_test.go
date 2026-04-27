package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"transaction-core-service/internal/controllers"
	"transaction-core-service/internal/models"
	"transaction-core-service/internal/repositories"
	"transaction-core-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTransactionRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	repo := repositories.NewTransactionRepository()
	svc := services.NewTransactionService(repo)
	ctrl := controllers.NewTransactionController(svc)

	router := gin.New()
	router.POST("/transaction/initiate", ctrl.InitiateTransaction)
	router.POST("/api/transaction/:transaction_id/submit", ctrl.SubmitTransaction)
	router.POST("/api/transaction/:transaction_id/validate", ctrl.ValidateTransaction)
	router.POST("/api/transaction/:transaction_id/broadcast", ctrl.BroadcastTransaction)
	router.POST("/api/transaction/:transaction_id/confirm", ctrl.ConfirmTransaction)
	router.POST("/api/transaction/:transaction_id/cancel", ctrl.CancelTransaction)
	router.GET("/api/transaction/:transaction_id/status", ctrl.GetStatus)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	return router
}

func TestInitiateTransaction(t *testing.T) {
	// В тестах пропускаем RabbitMQ, установив переменную окружения
	// (например, можно вызвать os.Setenv("SKIP_RABBITMQ", "true") в начале теста)
	t.Setenv("SKIP_RABBITMQ", "true")

	router := setupTransactionRouter()

	tx := models.Transaction{
		TransactionID: "tx123",
		StoreID:       "store123",
		Timestamp:     time.Now(),
		Amount:        "1000000000000000000",
		// Currency не задаётся, т.к. контроллер установит его в "ETH"
	}
	jsonValue, err := json.Marshal(tx)
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/transaction/initiate", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "transaction initiated", resp["status"])
}

func TestTransactionLifecycleBroadcast(t *testing.T) {
	t.Setenv("SKIP_RABBITMQ", "true")

	router := setupTransactionRouter()
	tx := models.Transaction{
		TransactionID: "tx-lifecycle",
		StoreID:       "store123",
		Timestamp:     time.Now(),
		Amount:        "1000000000000000000",
	}
	jsonValue, err := json.Marshal(tx)
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/transaction/initiate", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	for _, path := range []string{
		"/api/transaction/tx-lifecycle/submit",
		"/api/transaction/tx-lifecycle/validate",
		"/api/transaction/tx-lifecycle/broadcast",
		"/api/transaction/tx-lifecycle/confirm",
	} {
		req, _ = http.NewRequest("POST", path, nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	req, _ = http.NewRequest("GET", "/api/transaction/tx-lifecycle/status", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Data models.TransactionStatusResponse `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "confirmed", resp.Data.Status)
	assert.NotEmpty(t, resp.Data.BlockchainTxHash)
}

func TestTransactionLifecycleRejectsInvalidTransitions(t *testing.T) {
	t.Setenv("SKIP_RABBITMQ", "true")

	router := setupTransactionRouter()
	tx := models.Transaction{
		TransactionID: "tx-invalid-flow",
		StoreID:       "store123",
		Timestamp:     time.Now(),
		Amount:        "1000000000000000000",
	}
	jsonValue, err := json.Marshal(tx)
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/transaction/initiate", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("POST", "/api/transaction/tx-invalid-flow/validate", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)

	for _, path := range []string{
		"/api/transaction/tx-invalid-flow/submit",
		"/api/transaction/tx-invalid-flow/validate",
	} {
		req, _ = http.NewRequest("POST", path, nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	req, _ = http.NewRequest("POST", "/api/transaction/tx-invalid-flow/cancel", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)

	for _, path := range []string{
		"/api/transaction/tx-invalid-flow/broadcast",
		"/api/transaction/tx-invalid-flow/confirm",
	} {
		req, _ = http.NewRequest("POST", path, nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	req, _ = http.NewRequest("POST", "/api/transaction/tx-invalid-flow/submit", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)
}
