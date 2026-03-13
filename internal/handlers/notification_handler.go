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
// @Description  Get list of notifications for the authenticated employee
// @Tags         notifications
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   models.Notification
// @Failure      401  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /notificaciones [get]
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	employeeID, ok := requireEmployeeID(c)
	if !ok {
		return
	}
	resp, err := h.service.GetNotifications(employeeID)
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
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /notificaciones/{id}/leida [patch]
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	employeeID, ok := requireEmployeeID(c)
	if !ok {
		return
	}
	id := c.Param("id")
	if err := h.service.MarkAsRead(employeeID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// MarkAllAsRead godoc
// @Summary      Mark all notifications as read
// @Description  Mark all notifications for the authenticated employee as read
// @Tags         notifications
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  models.MarkAllReadResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /notificaciones/leer-todas [patch]
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	employeeID, ok := requireEmployeeID(c)
	if !ok {
		return
	}
	count, err := h.service.MarkAllAsRead(employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "count": count})
}
