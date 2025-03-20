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
