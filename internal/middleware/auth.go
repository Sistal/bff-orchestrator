package middleware

import (
	"net/http"
	"strings"

	"github.com/Sistal/bff-orchestrator/internal/services"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Authorization header required"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Invalid authorization format"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Call real auth service to validate
		resp, err := authService.Validate(token)
		if err != nil {
			// If services are down, we might want to return 500 or 401. 
			// Assuming validation error means unauthorized here.
			// But if it's connection error, it's different. 
			// For simplicity:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Token validation failed"})
			c.Abort()
			return
		}

		if !resp.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Invalid token"})
			c.Abort()
			return
		}
		
		// Set context variables
		// Cast to string if needed, or keep as int. Most handlers expect string conversion or interface{}
		c.Set("userID", resp.UserID)
		c.Set("username", resp.Username)
		c.Set("role", resp.Role)

		c.Next()
	}
}
