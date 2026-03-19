package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// requireEmployeeID extrae el employeeID del contexto y retorna false (con 401) si está vacío.
// Usado por handlers que operan sobre el funcionario autenticado.
func requireEmployeeID(c *gin.Context) (string, bool) {
	// La clave "employeeID" es establecida por el middleware BearerAuth
	id := c.GetString("employeeID")
	if id == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Usuario sin funcionario asociado",
		})
		return "", false
	}
	return id, true
}
