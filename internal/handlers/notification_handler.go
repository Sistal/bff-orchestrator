package handlers

import (
	"net/http"

	"github.com/Sistal/bff-orchestrator/internal/services"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service services.NotificationService
}

func NewNotificationHandler(s services.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: s}
}

// GetNotifications godoc
// @Summary      Get notifications
// @Description  Get list of notifications for the current user
// @Tags         notifications
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /notificaciones [get]
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID := c.GetString("userID")
	resp, err := h.service.GetNotifications(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// MarkAsRead godoc
// @Summary      Mark notification as read
// @Description  Mark a specific notification as read
// @Tags         notifications
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Notification ID"
// @Success      200
// @Failure      500  {object}  map[string]string
// @Router       /notificaciones/{id}/leida [patch]
// MarkAllAsRead godoc
// @Summary      Mark all notifications as read
// @Description  Mark all notifications for the current user as read
// @Tags         notifications
// @Security     BearerAuth
// @Produce      json
// @Success      200
// @Failure      500  {object}  map[string]string
// @Router       /notificaciones/leer-todas [patch]
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetString("userID")
	id := c.Param("id")
	if err := h.service.MarkAsRead(userID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetString("userID")
	if err := h.service.MarkAllAsRead(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
