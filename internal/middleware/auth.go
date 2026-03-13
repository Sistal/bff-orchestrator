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

// BearerAuthMiddleware valida el JWT Bearer Token, propaga todos los claims del token
// en el contexto de Gin y resuelve id_usuario → id_funcionario consultando al ms-funcionario.
//
// Claves disponibles en el contexto después de este middleware:
//   - "userID"        — id_usuario (string)
//   - "username"      — nombre_usuario (string)
//   - "nombreCompleto"— nombre_completo (string)
//   - "rut"           — rut del usuario (string)
//   - "rolID"         — id_rol (string)
//   - "rolName"       — nombre_rol (string)
//   - "estadoID"      — id_estado_usuario (string)
//   - "estadoName"    — nombre_estado (string)
//   - "employeeID"    — id_funcionario (string); vacío si el usuario no tiene funcionario asociado
//
// DEUDA TÉCNICA: la resolución userID→employeeID realiza una llamada HTTP extra por request.
// Implementar caché con TTL en el futuro para reducir latencia.
func BearerAuthMiddleware(identityClient *clients.IdentityClient, hrClient *clients.HRClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.Get()
		ip := c.ClientIP()
		path := c.FullPath()

		// ── 1. Extraer token desde header Authorization o cookie ──────────────
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

		// ── 2. Validar token contra ms-authentication ─────────────────────────
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

		// ── 3. Propagar todos los claims del token en el contexto ─────────────
		userIDStr := fmt.Sprintf("%d", resp.UserID)

		c.Set("userID", userIDStr)
		c.Set("username", resp.Username)
		c.Set("nombreCompleto", resp.NombreCompleto)
		c.Set("rut", resp.RUT)
		c.Set("rolID", fmt.Sprintf("%d", resp.Role))
		c.Set("rolName", resp.RoleName)
		c.Set("estadoID", fmt.Sprintf("%d", resp.EstadoID))
		c.Set("estadoName", resp.EstadoName)

		log.Debug("Auth: token validado correctamente",
			zap.String("user_id", userIDStr),
			zap.String("username", resp.Username),
			zap.String("role", resp.RoleName),
			zap.String("path", path),
			zap.String("ip", ip),
		)

		// ── 4. Resolver id_usuario → id_funcionario via ms-funcionario ────────
		// DEUDA TÉCNICA: esta llamada ocurre en cada request autenticado.
		// Implementar caché con TTL en el futuro para reducir latencia.
		employeeID, err := hrClient.GetEmployeeByUserID(userIDStr)
		if err != nil {
			// No abortamos: usuarios admin pueden no tener funcionario asociado.
			// Cada handler que requiera employeeID verificará si está vacío.
			log.Warn("Auth: no se pudo resolver id_funcionario para el usuario",
				zap.String("user_id", userIDStr),
				zap.String("path", path),
				zap.String("ip", ip),
				zap.Error(err),
			)
			c.Set("employeeID", "")
		} else if employeeID == 0 {
			log.Debug("Auth: usuario sin funcionario asociado (posible admin)",
				zap.String("user_id", userIDStr),
				zap.String("path", path),
			)
			c.Set("employeeID", "")
		} else {
			employeeIDStr := fmt.Sprintf("%d", employeeID)
			c.Set("employeeID", employeeIDStr)
			log.Debug("Auth: id_funcionario resuelto correctamente",
				zap.String("user_id", userIDStr),
				zap.String("employee_id", employeeIDStr),
				zap.String("path", path),
			)
		}

		c.Next()
	}
}
