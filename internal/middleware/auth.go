package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Sistal/bff-orchestrator/internal/clients"
	"github.com/Sistal/bff-orchestrator/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// APIKeyMiddleware valida el header x-api-key para comunicación server-to-server.
func APIKeyMiddleware(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("x-api-key")
		if key == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "x-api-key header required"})
			c.Abort()
			return
		}

		if key != apiKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Invalid API key"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// BearerAuthMiddleware valida el JWT Bearer Token y propaga userID en el contexto de Gin.
// Primero intenta leer el token del header Authorization; si está ausente, lo lee
// desde la cookie HttpOnly "access_token" establecida por el handler de Login.
// Después de este middleware, los handlers pueden usar c.GetString("userID").
func BearerAuthMiddleware(identityClient *clients.IdentityClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.Get()
		ip := c.ClientIP()
		path := c.FullPath()

		var token string

		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			cookieToken, err := c.Cookie("access_token")
			if err != nil || cookieToken == "" {
				log.Warn("Auth: token ausente (header y cookie)",
					zap.String("path", path),
					zap.String("ip", ip),
				)
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":   "Unauthorized",
					"message": "Authorization Bearer token required",
				})
				c.Abort()
				return
			}
			token = cookieToken
		}
		resp, err := identityClient.ValidateToken(token)
		if err != nil {
			log.Error("Auth: error al validar token contra ms-authentication",
				zap.String("path", path),
				zap.String("ip", ip),
				zap.Error(err),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		if !resp.Valid {
			log.Warn("Auth: token inválido o expirado",
				zap.String("path", path),
				zap.String("ip", ip),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		userIDStr := fmt.Sprintf("%d", resp.UserID)
		log.Debug("Auth: token validado correctamente",
			zap.String("user_id", userIDStr),
			zap.String("username", resp.Username),
			zap.String("role", resp.RoleName),
			zap.String("path", path),
			zap.String("ip", ip),
		)

		// Propagar userID como string para c.GetString("userID") en handlers
		c.Set("userID", userIDStr)
		c.Next()
	}
}
