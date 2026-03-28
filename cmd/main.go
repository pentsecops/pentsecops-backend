package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pentsecops/backend/internal/adapter/db/postgres"
	"github.com/pentsecops/backend/internal/adapter/http/handlers"
	"github.com/pentsecops/backend/internal/adapter/http/routes"
	"github.com/pentsecops/backend/internal/config"
	"github.com/pentsecops/backend/internal/core/usecases"
	"github.com/pentsecops/backend/internal/infra/cache"
	"github.com/pentsecops/backend/internal/infra/database"
	"github.com/pentsecops/backend/pkg/auth"
	"github.com/pentsecops/backend/pkg/auth/logger"
	"github.com/pentsecops/backend/pkg/middleware"
)

func main() {
	// Initialize logger
	appLogger := logger.GetLogger()
	logFileName := fmt.Sprintf("logs/pentsecops_%s.log", time.Now().Format("2006-01-02"))

	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}

	if err := appLogger.Init(logFileName); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer appLogger.Close()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load configuration: %v", err)
		log.Fatalf("Failed to load configuration: %v", err)
	}

	logger.Info("Starting PentSecOps Backend - Env: %s, Port: %s", cfg.Server.Env, cfg.Server.Port)

	// Initialize database
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to database: %v", err)
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	logger.Info("Database connected successfully")

	// Auto-initialize schema (create tables if they don't exist)
	if err := database.InitializeSchema(db); err != nil {
		logger.Error("Failed to initialize database schema: %v", err)
		log.Fatalf("Failed to initialize database schema: %v", err)
	}
	logger.Info("Database schema initialized successfully")

	// Initialize cache
	cache, err := cache.NewCache()
	if err != nil {
		logger.Error("Failed to initialize cache: %v", err)
		log.Fatalf("Failed to initialize cache: %v", err)
	}
	defer cache.Close()
	logger.Info("Cache initialized successfully")

	// Initialize PASETO service
	pasetoService, err := auth.NewPasetoService(
		cfg.Auth.PrivateKey,
		cfg.Auth.PublicKey,
	)
	if err != nil {
		logger.Error("Failed to initialize PASETO service: %v", err)
		log.Fatalf("Failed to initialize PASETO service: %v", err)
	}
	logger.Info("PASETO service initialized successfully")

	// Initialize repositories
	adminOverviewRepo := postgres.NewAdminOverviewRepository(db)
	usersRepo := postgres.NewUsersRepository(db)
	projectsRepo := postgres.NewProjectsRepository(db)
	tasksRepo := postgres.NewTasksRepository(db)
	vulnerabilitiesRepo := postgres.NewVulnerabilitiesRepository(db)
	domainsRepo := postgres.NewDomainsRepository(db)
	notificationsRepo := postgres.NewNotificationsRepository(db)
	auditRepo := postgres.NewAuditRepository(db)
	authRepo := postgres.NewAuthRepository(db)

	// Initialize use cases
	adminOverviewUseCase := usecases.NewAdminOverviewUseCase(adminOverviewRepo)
	usersUseCase := usecases.NewUsersUseCase(usersRepo)
	cachedUsersUseCase := usecases.NewCachedUsersUseCase(usersUseCase, cache)
	projectsUseCase := usecases.NewProjectsUseCase(projectsRepo)
	cachedProjectsUseCase := usecases.NewCachedProjectsUseCase(projectsUseCase, cache)
	tasksUseCase := usecases.NewTasksUseCase(tasksRepo)
	vulnerabilitiesUseCase := usecases.NewVulnerabilitiesUseCase(vulnerabilitiesRepo)
	cachedVulnerabilitiesUseCase := usecases.NewCachedVulnerabilitiesUseCase(vulnerabilitiesUseCase, cache.GetClient())
	domainsUseCase := usecases.NewDomainsUseCase(domainsRepo)
	cachedDomainsUseCase := usecases.NewCachedDomainsUseCase(domainsUseCase, cache.GetClient())
	notificationsUseCase := usecases.NewNotificationsUseCase(notificationsRepo)
	cachedNotificationsUseCase := usecases.NewCachedNotificationsUseCase(notificationsUseCase, cache.GetClient())
	adminNotificationsUseCase := usecases.NewAdminNotificationsUseCase(notificationsRepo, usersRepo)
	authUseCase := usecases.NewAuthUseCase(authRepo, pasetoService, cfg.Auth.AccessTokenDuration, cfg.Auth.RefreshTokenDuration)
	pentesterOverviewRepo := postgres.NewPentesterOverviewRepository(db)
	pentesterOverviewUseCase := usecases.NewPentesterOverviewUseCase(pentesterOverviewRepo)
	pentesterProjectsRepo := postgres.NewPentesterProjectsRepository(db)
	pentesterProjectsUseCase := usecases.NewPentesterProjectsUseCase(pentesterProjectsRepo)
	pentesterTasksRepo := postgres.NewPentesterTasksRepository(db)
	pentesterTasksUseCase := usecases.NewPentesterTasksUseCase(pentesterTasksRepo)
	pentesterSubmitReportRepo := postgres.NewPentesterSubmitReportRepository(db)
	pentesterSubmitReportUseCase := usecases.NewPentesterSubmitReportUseCase(pentesterSubmitReportRepo)
	pentesterAlertsRepo := postgres.NewPentesterAlertsRepository(db)
	pentesterAlertsUseCase := usecases.NewPentesterAlertsUseCase(pentesterAlertsRepo)
	stakeholderOverviewRepo := postgres.NewStakeholderOverviewRepository(db)
	stakeholderOverviewUseCase := usecases.NewStakeholderOverviewUseCase(stakeholderOverviewRepo)
	stakeholderVulnerabilitiesRepo := postgres.NewStakeholderVulnerabilitiesRepository(db)
	stakeholderVulnerabilitiesUseCase := usecases.NewStakeholderVulnerabilitiesUseCase(stakeholderVulnerabilitiesRepo)
	stakeholderReportsRepo := postgres.NewStakeholderReportsRepository(db)
	stakeholderReportsUseCase := usecases.NewStakeholderReportsUseCase(stakeholderReportsRepo)

	// Initialize handlers
	adminOverviewHandler := handlers.NewAdminOverviewHandler(adminOverviewUseCase)
	usersHandler := handlers.NewUsersHandler(cachedUsersUseCase)
	projectsHandler := handlers.NewProjectsHandler(cachedProjectsUseCase)
	tasksHandler := handlers.NewTasksHandler(tasksUseCase)
	vulnerabilitiesHandler := handlers.NewVulnerabilitiesHandler(cachedVulnerabilitiesUseCase)
	domainsHandler := handlers.NewDomainsHandler(cachedDomainsUseCase)
	notificationsHandler := handlers.NewNotificationsHandler(cachedNotificationsUseCase)
	auditHandler := handlers.NewAuditHandler(auditRepo)
	adminNotificationsHandler := handlers.NewAdminNotificationsHandler(adminNotificationsUseCase)
	authHandler := handlers.NewAuthHandler(authUseCase, appLogger)
	pentesterOverviewHandler := handlers.NewPentesterOverviewHandler(pentesterOverviewUseCase)
	pentesterProjectsHandler := handlers.NewPentesterProjectsHandler(pentesterProjectsUseCase)
	pentesterTasksHandler := handlers.NewPentesterTasksHandler(pentesterTasksUseCase)
	pentesterSubmitReportHandler := handlers.NewPentesterSubmitReportHandler(pentesterSubmitReportUseCase)
	pentesterAlertsHandler := handlers.NewPentesterAlertsHandler(pentesterAlertsUseCase)
	stakeholderOverviewHandler := handlers.NewStakeholderOverviewHandler(stakeholderOverviewUseCase)
	stakeholderVulnerabilitiesHandler := handlers.NewStakeholderVulnerabilitiesHandler(stakeholderVulnerabilitiesUseCase)
	stakeholderReportsHandler := handlers.NewStakeholderReportsHandler(stakeholderReportsUseCase)

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(pasetoService)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:               "PentSecOps API",
		ServerHeader:          "PentSecOps",
		DisableStartupMessage: false,
		ErrorHandler:          customErrorHandler,
		Prefork:               false,
		StrictRouting:         false,
		CaseSensitive:         false,
		BodyLimit:             10 * 1024 * 1024, // 10MB
	})

	// Apply global middleware
	app.Use(middleware.Recovery())
	app.Use(middleware.CORS(&cfg.CORS))
	app.Use(middleware.RateLimit(&cfg.RateLimit))

	// Setup routes
	fmt.Printf("=== SETTING UP ROUTES ===\n")
	llmHandler := handlers.NewLLMHandler()
	routes.SetupRoutes(app, authHandler, authMiddleware, adminOverviewHandler, usersHandler, projectsHandler, tasksHandler, vulnerabilitiesHandler, domainsHandler, notificationsHandler, auditHandler, adminNotificationsHandler, pentesterOverviewHandler, pentesterProjectsHandler, pentesterTasksHandler, pentesterSubmitReportHandler, pentesterAlertsHandler, stakeholderOverviewHandler, stakeholderVulnerabilitiesHandler, stakeholderReportsHandler, llmHandler)
	fmt.Printf("=== ROUTES SETUP COMPLETE ===\n")

	// Start server in a goroutine
	go func() {
		// Prefer PORT env var (used by Railway), fallback to cfg.Server.Port
		port := os.Getenv("PORT")
		if port == "" {
			port = cfg.Server.Port
		}
		addr := fmt.Sprintf("0.0.0.0:%s", port)
		logger.Info("Server starting on %s", addr)
		fmt.Printf("=== SERVER STARTING ON %s ===\n", addr)

		if err := app.Listen(addr); err != nil {
			logger.Error("Failed to start server: %v", err)
			fmt.Printf("=== SERVER START ERROR: %v ===\n", err)
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	logger.Info("PentSecOps Backend is running. Press Ctrl+C to stop.")

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		logger.Error("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited successfully")
}

// customErrorHandler handles errors globally
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	log.Printf("Request error: %v, path: %s, method: %s, status: %d", err, c.Path(), c.Method(), code)

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    "ERROR",
			"message": err.Error(),
		},
	})
}
