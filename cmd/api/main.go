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

	_ "github.com/Sistal/bff-orchestrator/docs"
	"github.com/Sistal/bff-orchestrator/internal/clients"
	"github.com/Sistal/bff-orchestrator/internal/config"
	"github.com/Sistal/bff-orchestrator/internal/handlers"
	"github.com/Sistal/bff-orchestrator/internal/middleware"
	"github.com/Sistal/bff-orchestrator/internal/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           BFF Orchestrator API
// @version         1.0
// @description     BFF Orchestrator API for Sistal Application.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Load configuration
	cfg := config.Load()

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Clients
	identityClient := clients.NewIdentityClient()
	hrClient := clients.NewHRClient()
	opsClient := clients.NewOpsClient()
	catalogClient := clients.NewCatalogClient()

	// Initialize Services (Real HTTP implementations)
	authService := services.NewHTTPAuthService(identityClient)
	catalogService := services.NewHTTPCatalogService(catalogClient)
	branchService := services.NewHTTPBranchService(hrClient)

	// Ops Service handles Request, Delivery and Notifications
	opsService := services.NewHTTPOpsService(opsClient)

	// Employee Service combines HR and Ops (stats)
	employeeService := services.NewHTTPEmployeeService(hrClient, opsClient)

	// Cast OpsService to interfaces required by handlers
	var requestService services.RequestService = opsService
	var deliveryService services.DeliveryService = opsService
	var notificationService services.NotificationService = opsService

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	catalogHandler := handlers.NewCatalogHandler(catalogService)
	branchHandler := handlers.NewBranchHandler(branchService)
	requestHandler := handlers.NewRequestHandler(requestService)
	deliveryHandler := handlers.NewDeliveryHandler(deliveryService)
	employeeHandler := handlers.NewEmployeeHandler(employeeService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	healthHandler := handlers.NewHealthHandler()

	// Initialize Router
	r := gin.Default()
	r.Use(middleware.CORS())

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public Routes
	r.GET("/health", healthHandler.Check)
	r.GET("/auth/validate", authHandler.Validate)

	// Protected Routes
	api := r.Group("/")
	api.Use(middleware.AuthMiddleware(authService))
	{
		// Auth
		api.GET("/auth/me", authHandler.GetMe)

		// Catalog
		api.GET("/catalogo/tallas", catalogHandler.GetSizes)
		api.GET("/catalogo/motivos-cambio", catalogHandler.GetChangeReasons)
		api.GET("/catalogo/prenda-tipos", catalogHandler.GetGarmentTypes)
		api.GET("/campanas/activa", catalogHandler.GetActiveCampaign)

		// Branches
		api.GET("/sucursales", branchHandler.GetAllBranches)
		api.GET("/solicitudes/cambio-sucursal/historial", branchHandler.GetChangeHistory)
		api.POST("/solicitudes/cambio-sucursal", branchHandler.CreateChangeRequest)

		// Requests
		api.GET("/solicitudes", requestHandler.GetRequests)
		api.GET("/solicitudes/:id", requestHandler.GetRequestByID)
		api.POST("/solicitudes/reposicion", requestHandler.CreateReplenishmentRequest)
		api.POST("/solicitudes/cambio-prenda", requestHandler.CreateGarmentChangeRequest)
		api.POST("/archivos/upload", requestHandler.UploadFile)

		// Deliveries
		api.GET("/entregas", deliveryHandler.GetDeliveries)
		api.GET("/entregas/:id", deliveryHandler.GetDeliveryByID)
		api.POST("/entregas/:id/confirmar", deliveryHandler.ConfirmDelivery)

		// Notifications
		api.GET("/notificaciones", notificationHandler.GetNotifications)
		api.PATCH("/notificaciones/:id/leida", notificationHandler.MarkAsRead)
		api.PATCH("/notificaciones/leer-todas", notificationHandler.MarkAllAsRead)

		// Employee / Officials
		v1 := api.Group("/api/v1")
		{
			funcionarios := v1.Group("/funcionarios")
			{
				funcionarios.GET("/me", employeeHandler.GetProfile)
				funcionarios.PUT("/me", employeeHandler.UpdateContact)
				funcionarios.PUT("/me/preferencias", employeeHandler.UpdatePreferences)
				funcionarios.PUT("/me/seguridad", employeeHandler.UpdateSecurity)
				funcionarios.GET("/me/stats", employeeHandler.GetStats)
				funcionarios.GET("/me/actividad", employeeHandler.GetActivity)
				funcionarios.GET("/:id/medidas", employeeHandler.GetMeasurements)
				funcionarios.POST("/:id/medidas", employeeHandler.RegisterMeasurements)
			}
		}
	}

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
