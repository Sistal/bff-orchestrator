package handlers

import (
	"net/http"

	"github.com/Sistal/bff-orchestrator/internal/services"
	"github.com/gin-gonic/gin"
)

type DeliveryHandler struct {
	service services.DeliveryService
}

func NewDeliveryHandler(s services.DeliveryService) *DeliveryHandler {
	return &DeliveryHandler{service: s}
}

// GetDeliveries godoc
// @Summary      Get deliveries
// @Description  Get list of deliveries for the current user
// @Tags         deliveries
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /entregas [get]
func (h *DeliveryHandler) GetDeliveries(c *gin.Context) {
	userID := c.GetString("userID")
	resp, err := h.service.GetDeliveries(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetDeliveryByID godoc
// @Summary      Get delivery by ID
// @Description  Get a specific delivery by its ID
// @Tags         deliveries
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Delivery ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]string
// @Router       /entregas/{id} [get]
func (h *DeliveryHandler) GetDeliveryByID(c *gin.Context) {
	userID := c.GetString("userID")
	id := c.Param("id")
	resp, err := h.service.GetDeliveryByID(userID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "Delivery not found"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// ConfirmDelivery godoc
// @Summary      Confirm delivery
// @Description  Confirm receipt of a delivery
// @Tags         deliveries
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Delivery ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /entregas/{id}/confirmar [post]
func (h *DeliveryHandler) ConfirmDelivery(c *gin.Context) {
	userID := c.GetString("userID")
	id := c.Param("id")
	resp, err := h.service.ConfirmDelivery(userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
