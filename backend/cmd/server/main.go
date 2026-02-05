package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/udai-kiran/agentic-cash/internal/application/service"
	"github.com/udai-kiran/agentic-cash/internal/config"
	"github.com/udai-kiran/agentic-cash/internal/infrastructure/auth"
	"github.com/udai-kiran/agentic-cash/internal/infrastructure/persistence/postgres"
	httpRouter "github.com/udai-kiran/agentic-cash/internal/interfaces/http"
	"github.com/udai-kiran/agentic-cash/internal/interfaces/http/handler"
	"github.com/udai-kiran/agentic-cash/pkg/logger"
)

func main() {
	// Initialize logger
	logger.Init(os.Getenv("GO_ENV") == "production")

	// Load configuration
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// Create context for database initialization
	ctx := context.Background()

	// Initialize database connection pool
	pool, err := postgres.NewPool(ctx, &cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	logger.Info("Connected to PostgreSQL successfully")

	// Initialize application tables
	if err := postgres.InitializeAppTables(ctx, pool); err != nil {
		logger.Error("Failed to initialize app tables", "error", err)
		os.Exit(1)
	}

	logger.Info("Application tables initialized")

	// Start token cleanup service (runs every 6 hours)
	tokenCleanup := postgres.NewTokenCleanupService(pool, 6*time.Hour)
	go tokenCleanup.Start(ctx)
	defer tokenCleanup.Stop()

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
	)

	// Initialize repositories
	accountRepo := postgres.NewAccountRepository(pool)
	userRepo := postgres.NewUserRepository(pool)
	transactionRepo := postgres.NewTransactionRepository(pool)
	commodityRepo := postgres.NewCommodityRepository(pool)

	// Initialize services
	authService := service.NewAuthService(userRepo, jwtManager)
	analyticsService := service.NewAnalyticsService(accountRepo, transactionRepo)

	// Initialize handlers
	accountHandler := handler.NewAccountHandler(accountRepo, commodityRepo)
	authHandler := handler.NewAuthHandler(authService)
	transactionHandler := handler.NewTransactionHandler(transactionRepo)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)
	commodityHandler := handler.NewCommodityHandler(commodityRepo)

	// Setup router
	router := httpRouter.Router(&httpRouter.RouterConfig{
		AccountHandler:     accountHandler,
		AuthHandler:        authHandler,
		TransactionHandler: transactionHandler,
		AnalyticsHandler:   analyticsHandler,
		CommodityHandler:   commodityHandler,
		JWTManager:         jwtManager,
		AllowedOrigins:     cfg.CORS.AllowedOrigins,
	})

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting server", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("Server exited")
}
