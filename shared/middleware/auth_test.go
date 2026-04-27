package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"shared/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuthMiddlewareRejectsMissingHeader(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(JWTAuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestJWTAuthMiddlewareAllowsValidToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")
	gin.SetMode(gin.TestMode)

	token, err := utils.GenerateTokenWithRoles("1", []string{"admin"}, time.Hour)
	assert.NoError(t, err)

	router := gin.New()
	router.Use(JWTAuthMiddleware(), RoleMiddleware("admin"))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

