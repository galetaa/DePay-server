package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func DemoEndpointGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		appEnv := strings.ToLower(os.Getenv("APP_ENV"))
		demoEnabled := strings.EqualFold(os.Getenv("DEMO_ENDPOINTS_ENABLED"), "true")
		if appEnv == "production" && !demoEnabled {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "demo endpoints are disabled in production"})
			return
		}
		c.Next()
	}
}

func ProductionRoleGate(requiredRoles ...string) gin.HandlerFunc {
	authMiddleware := JWTAuthMiddleware()
	roleMiddleware := RoleMiddleware(requiredRoles...)

	return func(c *gin.Context) {
		if strings.ToLower(os.Getenv("APP_ENV")) != "production" && !strings.EqualFold(os.Getenv("REQUIRE_AUTH"), "true") {
			c.Next()
			return
		}
		authMiddleware(c)
		if c.IsAborted() {
			return
		}
		roleMiddleware(c)
	}
}
