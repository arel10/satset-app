package middleware

import (
	"github.com/gin-gonic/gin"

	"satset-app/internal/handler"
)

// RequireRole checks whether the authenticated user has one of the required roles.
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			handler.RespondUnauthorized(c, "Autentikasi diperlukan")
			c.Abort()
			return
		}

		userRole, ok := roleVal.(string)
		if !ok {
			handler.RespondUnauthorized(c, "Role pengguna tidak valid")
			c.Abort()
			return
		}

		for _, allowed := range allowedRoles {
			if userRole == allowed {
				c.Next()
				return
			}
		}

		handler.RespondForbidden(c, "Anda tidak memiliki hak akses untuk tindakan ini")
		c.Abort()
	}
}
