package middleware

import (
	"os"
	"strings"

	"github.com/gin-contrib/cors"
)

// CORSConfig construye la configuración CORS a partir de la variable de entorno
// FRONTEND_ORIGINS (lista separada por comas) o FRONTEND_ORIGIN.
// Si el valor es "*" o está vacío, se permite cualquier origen (solo desarrollo).
//
// Uso en main:
//
//	r.Use(cors.New(middleware.CORSConfig()))
func CORSConfig() cors.Config {
	raw := strings.TrimSpace(os.Getenv("FRONTEND_ORIGINS"))
	if raw == "" {
		raw = strings.TrimSpace(os.Getenv("FRONTEND_ORIGIN"))
	}

	cfg := cors.DefaultConfig()
	cfg.AllowCredentials = true
	cfg.AllowHeaders = []string{
		"Origin", "Content-Type", "Content-Length",
		"Accept-Encoding", "X-CSRF-Token", "Authorization",
		"Accept", "Cache-Control", "X-Requested-With",
	}
	cfg.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}

	if raw == "" || raw == "*" {
		// AllowAllOrigins=true envía "*", incompatible con AllowCredentials.
		// AllowOriginFunc refleja el origen exacto → compatible con cookies/credenciales.
		cfg.AllowOriginFunc = func(origin string) bool { return true }
	} else {
		origins := strings.Split(raw, ",")
		for i, o := range origins {
			origins[i] = strings.TrimSpace(o)
		}
		cfg.AllowOrigins = origins
	}

	return cfg
}
