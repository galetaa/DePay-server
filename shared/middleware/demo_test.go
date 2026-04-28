package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestDemoEndpointGuardRejectsProductionByDefault(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("DEMO_ENDPOINTS_ENABLED", "")
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/demo", DemoEndpointGuard(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodPost, "/demo", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestDemoEndpointGuardAllowsLocal(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/demo", DemoEndpointGuard(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodPost, "/demo", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
