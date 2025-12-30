package handler

import (
	"net/http"
	"strconv"

	"github.com/Sistal/bff-orchestrator/internal/ports"
	"github.com/gin-gonic/gin"
)

// AggregationHandler handles aggregation requests
type AggregationHandler struct {
	aggregationService ports.AggregationService
}

// NewAggregationHandler creates a new aggregation handler
func NewAggregationHandler(aggregationService ports.AggregationService) *AggregationHandler {
	return &AggregationHandler{
		aggregationService: aggregationService,
	}
}

// GetDashboard returns aggregated dashboard data
func (h *AggregationHandler) GetDashboard(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	dashboard, err := h.aggregationService.GetDashboard(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dashboard)
}
