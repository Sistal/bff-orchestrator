package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sistal/bff-orchestrator/internal/adapters/http/handler"
	"github.com/Sistal/bff-orchestrator/internal/adapters/http/router"
	"github.com/Sistal/bff-orchestrator/internal/adapters/provider"
	"github.com/Sistal/bff-orchestrator/internal/application/service"
	"github.com/Sistal/bff-orchestrator/internal/config"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize external API providers (adapters)
	userProvider := provider.NewUserProvider(cfg.UserServiceURL)
	productProvider := provider.NewProductProvider(cfg.ProductServiceURL)

	// Initialize application services (use cases)
	aggregationService := service.NewAggregationService(userProvider, productProvider)

	// Initialize HTTP handlers
	healthHandler := handler.NewHealthHandler()
	aggregationHandler := handler.NewAggregationHandler(aggregationService)

	// Setup router
	r := router.SetupRouter(healthHandler, aggregationHandler)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting BFF Orchestrator on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
