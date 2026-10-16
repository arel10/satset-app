package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"satset-app/internal/handler"
	"satset-app/internal/service"
)

// AuthMiddleware extracts and validates Bearer token from the Authorization header.
func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			handler.RespondUnauthorized(c, "Token autentikasi tidak ditemukan")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			handler.RespondUnauthorized(c, "Format token autentikasi tidak valid")
			c.Abort()
			return
		}

		tokenStr := parts[1]
		claims, err := authService.ValidateAccessToken(tokenStr)
		if err != nil {
			handler.RespondUnauthorized(c, "Sesi tidak valid atau telah berakhir")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}
