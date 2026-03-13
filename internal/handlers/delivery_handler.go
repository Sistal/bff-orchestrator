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
// @Description  Get list of deliveries for the authenticated employee
// @Tags         deliveries
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /entregas [get]
func (h *DeliveryHandler) GetDeliveries(c *gin.Context) {
	employeeID, ok := requireEmployeeID(c)
	if !ok {
		return
	}
	resp, err := h.service.GetDeliveries(employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetDeliveryByID godoc
// @Summary      Get delivery by ID
// @Description  Get a specific delivery by its ID for the authenticated employee
// @Tags         deliveries
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Delivery ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /entregas/{id} [get]
func (h *DeliveryHandler) GetDeliveryByID(c *gin.Context) {
	employeeID, ok := requireEmployeeID(c)
	if !ok {
		return
	}
	id := c.Param("id")
	resp, err := h.service.GetDeliveryByID(employeeID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "Delivery not found"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// ConfirmDelivery godoc
// @Summary      Confirm delivery
// @Description  Confirm receipt of a delivery for the authenticated employee
// @Tags         deliveries
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Delivery ID"
// @Success      200  {object}  models.DeliverySummary
// @Failure      401  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /entregas/{id}/confirmar [post]
func (h *DeliveryHandler) ConfirmDelivery(c *gin.Context) {
	employeeID, ok := requireEmployeeID(c)
	if !ok {
		return
	}
	id := c.Param("id")
	resp, err := h.service.ConfirmDelivery(employeeID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetUpcomingDeliveries godoc
// @Summary      Get upcoming deliveries
// @Description  Get upcoming deliveries for the authenticated employee (dashboard)
// @Tags         deliveries
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   models.DeliverySummary
// @Failure      401  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /api/v1/entregas/upcoming [get]
func (h *DeliveryHandler) GetUpcomingDeliveries(c *gin.Context) {
	employeeID, ok := requireEmployeeID(c)
	if !ok {
		return
	}
	resp, err := h.service.GetUpcomingDeliveries(employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
