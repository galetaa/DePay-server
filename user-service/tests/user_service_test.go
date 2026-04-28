package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"user-service/internal/controllers"
	"user-service/internal/repositories"
	"user-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupRouter инициализирует Gin с базовыми эндпоинтами User Service
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	userRepo := repositories.NewUserRepository()
	userSvc := services.NewUserService(userRepo)
	userCtrl := controllers.NewUserController(userSvc)

	router := gin.New()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.POST("/user/register", userCtrl.Register)
	router.POST("/user/login", userCtrl.Login)
	router.POST("/user/refresh-token", userCtrl.RefreshToken)
	router.POST("/user/logout", userCtrl.Logout)

	return router
}

// TestHealthEndpoint проверяет работу эндпоинта /health
func TestHealthEndpoint(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}

// TestRegisterEndpoint проверяет регистрацию нового пользователя
func TestRegisterEndpoint(t *testing.T) {
	router := setupRouter()

	reqBody := map[string]string{
		"email":      "test@example.com",
		"first_name": "John",
		"last_name":  "Doe",
		"password":   "password123",
	}
	jsonValue, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/user/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.NotNil(t, resp["user"], "Response should contain a user object")
	assert.NotNil(t, resp["token"], "Response should contain a token")
}

func TestRefreshTokenRotatesAndRevokesOldToken(t *testing.T) {
	router := setupRouter()

	reqBody := map[string]string{
		"email":      "rotate@example.com",
		"first_name": "Rhea",
		"last_name":  "Token",
		"password":   "password123",
	}
	jsonValue, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/user/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var registerResp map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &registerResp))
	oldRefresh := registerResp["refresh_token"].(string)

	refreshBody, _ := json.Marshal(map[string]string{"token": oldRefresh})
	req, _ = http.NewRequest("POST", "/user/refresh-token", bytes.NewBuffer(refreshBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var refreshResp map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &refreshResp))
	assert.NotEmpty(t, refreshResp["token"])
	assert.NotEmpty(t, refreshResp["refresh_token"])
	assert.NotEqual(t, oldRefresh, refreshResp["refresh_token"])

	req, _ = http.NewRequest("POST", "/user/refresh-token", bytes.NewBuffer(refreshBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogoutRevokesRefreshToken(t *testing.T) {
	router := setupRouter()

	reqBody := map[string]string{
		"email":      "logout@example.com",
		"first_name": "Lora",
		"last_name":  "Out",
		"password":   "password123",
	}
	jsonValue, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/user/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var registerResp map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &registerResp))
	refresh := registerResp["refresh_token"].(string)

	body, _ := json.Marshal(map[string]string{"token": refresh})
	req, _ = http.NewRequest("POST", "/user/logout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	req, _ = http.NewRequest("POST", "/user/refresh-token", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
