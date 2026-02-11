package handlers

import (
	"net/http"
	"strings"

	"github.com/Sistal/bff-orchestrator/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler(s services.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

// Validate validates the token
// @Summary      Validate token
// @Description  Check if the token is valid
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer Token"
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]string
// @Router       /auth/validate [get]
func (h *AuthHandler) Validate(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	resp, err := h.service.Validate(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetMe returns current user info
// @Summary      Get current user
// @Description  Get currently logged in user details
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	resp, err := h.service.GetMe(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
