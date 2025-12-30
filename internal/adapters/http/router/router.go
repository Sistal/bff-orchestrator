package router

import (
	"github.com/Sistal/bff-orchestrator/internal/adapters/http/handler"
	"github.com/Sistal/bff-orchestrator/internal/adapters/http/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the Gin router
func SetupRouter(healthHandler *handler.HealthHandler, aggregationHandler *handler.AggregationHandler) *gin.Engine {
	r := gin.Default()

	// Apply CORS middleware
	r.Use(middleware.CORS())

	// Health check endpoint
	r.GET("/health", healthHandler.Check)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Dashboard endpoint - aggregates data from multiple sources
		v1.GET("/dashboard/:userId", aggregationHandler.GetDashboard)
	}

	return r
}
