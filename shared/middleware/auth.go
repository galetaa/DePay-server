package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"shared/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// JWTAuthMiddleware проверяет наличие и корректность JWT в заголовке запроса
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}
		tokenString := parts[1]

		// Здесь можно добавить проверку алгоритма, ключа и т.д.
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(utils.GetJWTSecret()), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// При необходимости можно сохранить claims в контексте
		c.Set("claims", token.Claims)
		c.Next()
	}
}

func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	required := make(map[string]struct{}, len(requiredRoles))
	for _, role := range requiredRoles {
		required[role] = struct{}{}
	}

	return func(c *gin.Context) {
		rawClaims, ok := c.Get("claims")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "claims are missing"})
			return
		}

		claims, ok := rawClaims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}

		rolesValue, ok := claims["roles"].([]interface{})
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "roles are missing"})
			return
		}

		for _, roleValue := range rolesValue {
			role, ok := roleValue.(string)
			if !ok {
				continue
			}
			if _, allowed := required[role]; allowed {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient role"})
	}
}
