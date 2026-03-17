package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/Sistal/bff-orchestrator/docs"
	"github.com/Sistal/bff-orchestrator/internal/clients"
	"github.com/Sistal/bff-orchestrator/internal/config"
	"github.com/Sistal/bff-orchestrator/internal/handlers"
	"github.com/Sistal/bff-orchestrator/internal/logger"
	"github.com/Sistal/bff-orchestrator/internal/middleware"
	"github.com/Sistal/bff-orchestrator/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
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
	// Cargar variables de entorno desde .env (si existe).
	// En producción con variables inyectadas por el orquestador, este paso es silencioso.
	if err := godotenv.Load(); err != nil {
		// No es fatal: en producción las variables ya están en el entorno del proceso.
		fmt.Println("INFO: archivo .env no encontrado, usando variables de entorno del sistema")
	}

	// Load configuration
	cfg := config.Load()

	// Inicializar logger (singleton)
	log := logger.Get()
	defer logger.Sync()

	log.Info("Iniciando BFF Orchestrator",
		zap.String("entorno", cfg.Environment),
		zap.String("puerto", cfg.Port),
	)

	// ── DIAGNÓSTICO: confirmar configuración efectiva de CORS y cookies ──────
	log.Info("Configuración efectiva de CORS y cookies",
		zap.String("frontend_origins_raw", os.Getenv("FRONTEND_ORIGINS")),
		zap.String("cookie_domain", cfg.CookieDomain),
		zap.Int("cookie_max_age", cfg.CookieMaxAge),
	)
	// ── FIN DIAGNÓSTICO ──────────────────────────────────────────────────────

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
	catalogService := services.NewHTTPCatalogService(catalogClient)
	branchService := services.NewHTTPBranchService(hrClient)

	// Ops Service handles Request, Delivery and Notifications
	opsService := services.NewHTTPOpsService(opsClient)

	// Employee Service combines HR and Ops (stats)
	employeeService := services.NewHTTPEmployeeService(hrClient, opsClient)

	// Auth Service now depends on Employee Service for full registration
	authService := services.NewHTTPAuthService(identityClient, employeeService)

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
	r.Use(cors.New(middleware.CORSConfig()))

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public Routes
	r.GET("/health", healthHandler.Check)

	// Auth — rutas públicas (sin JWT)
	authPublic := r.Group("/api/v1/auth")
	{
		authPublic.POST("/login", authHandler.Login)
		authPublic.POST("/register", authHandler.Register)
		authPublic.GET("/validate", authHandler.Validate)
		authPublic.POST("/refresh", authHandler.Refresh)
	}

	// Protected Routes — autenticadas con Bearer JWT
	api := r.Group("/")
	api.Use(middleware.BearerAuthMiddleware(identityClient, hrClient))
	{
		// Auth — rutas protegidas
		api.GET("/api/v1/auth/me", authHandler.GetMe)
		api.GET("/api/v1/auth/status", authHandler.Status)
		api.POST("/api/v1/auth/logout", authHandler.Logout)
		api.PUT("/api/v1/auth/change-password", authHandler.ChangePassword)
		api.GET("/api/v1/auth/roles", authHandler.GetRoles)

		// Auth — rutas de administración (requieren Admin o Super Admin)
		api.POST("/api/v1/auth/users", authHandler.CreateUser)
		api.GET("/api/v1/auth/users", authHandler.ListUsers)
		api.GET("/api/v1/auth/users/:id_usuario", authHandler.GetUserByID)
		api.PUT("/api/v1/auth/users/:id_usuario", authHandler.UpdateUser)

		// Catalog
		api.GET("/catalogo/tallas", catalogHandler.GetSizes)
		api.GET("/catalogo/motivos-cambio", catalogHandler.GetChangeReasons)
		api.GET("/catalogo/prenda-tipos", catalogHandler.GetGarmentTypes)
		api.GET("/campanas/activa", catalogHandler.GetActiveCampaign)

		// Branches
		api.GET("/sucursales", branchHandler.GetAllBranches)
		api.GET("/solicitudes/cambio-sucursal/historial", branchHandler.GetChangeHistory)
		api.POST("/solicitudes/cambio-sucursal", branchHandler.CreateChangeRequest)

		// Requests — rutas estáticas ANTES que /:id para evitar conflictos en Gin
		api.GET("/solicitudes", requestHandler.GetRequests)
		api.GET("/solicitudes/recent", requestHandler.GetRecentRequests)
		api.POST("/solicitudes/reposicion", requestHandler.CreateReplenishmentRequest)
		api.POST("/solicitudes/cambio-prenda", requestHandler.CreateGarmentChangeRequest)
		api.GET("/solicitudes/:id", requestHandler.GetRequestByID)
		api.POST("/archivos/upload", requestHandler.UploadFile)

		// Deliveries — rutas estáticas ANTES que /:id
		api.GET("/entregas", deliveryHandler.GetDeliveries)
		api.GET("/entregas/:id", deliveryHandler.GetDeliveryByID)
		api.POST("/entregas/:id/confirmar", deliveryHandler.ConfirmDelivery)

		// Notifications — rutas estáticas ANTES que /:id/leida
		api.GET("/notificaciones", notificationHandler.GetNotifications)
		api.PATCH("/notificaciones/leer-todas", notificationHandler.MarkAllAsRead)
		api.PATCH("/notificaciones/:id/leida", notificationHandler.MarkAsRead)
		// Employee / Officials
		v1 := api.Group("/api/v1")
		{
			// Entregas upcoming (dashboard)
			v1.GET("/entregas/upcoming", deliveryHandler.GetUpcomingDeliveries)

			funcionarios := v1.Group("/funcionarios")
			{
				// Rutas /me — deben ir ANTES que /:id
				funcionarios.GET("/me", employeeHandler.GetProfile)
				funcionarios.PUT("/me", employeeHandler.UpdateContact)
				funcionarios.PUT("/me/preferencias", employeeHandler.UpdatePreferences)
				funcionarios.PUT("/me/seguridad", employeeHandler.UpdateSecurity)
				funcionarios.GET("/me/stats", employeeHandler.GetStats)
				funcionarios.GET("/me/actividad", employeeHandler.GetActivity)

				// Ruta estática /filter — ANTES que /:id
				funcionarios.GET("/filter", employeeHandler.FilterEmployees)

				// CRUD Admin — 501 Not Implemented (pendiente de implementación en ms-funcionario)
				funcionarios.GET("", employeeHandler.ListEmployees)
				funcionarios.POST("", employeeHandler.CreateEmployee)

				// Rutas con :id — medidas/historial (estática) ANTES que medidas (dinámica)
				funcionarios.GET("/:id/medidas/historial", employeeHandler.GetMeasurementsHistory)
				funcionarios.GET("/:id/medidas", employeeHandler.GetMeasurements)
				funcionarios.POST("/:id/medidas", employeeHandler.RegisterMeasurements)
				funcionarios.PUT("/:id/medidas", employeeHandler.UpdateMeasurements)

				// CRUD Admin con :id — 501 Not Implemented
				funcionarios.GET("/:id", employeeHandler.GetEmployeeByID)
				funcionarios.PUT("/:id", employeeHandler.UpdateEmployee)
				funcionarios.DELETE("/:id", employeeHandler.DeleteEmployee)
				funcionarios.PATCH("/:id/activate", employeeHandler.ActivateEmployee)
				funcionarios.PATCH("/:id/deactivate", employeeHandler.DeactivateEmployee)
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
		log.Info("BFF Orchestrator escuchando", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Error fatal al iniciar el servidor", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	log.Info("Señal de apagado recibida, cerrando servidor...", zap.String("señal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Apagado forzado del servidor", zap.Error(err))
	}

	log.Info("Servidor apagado correctamente")
}
