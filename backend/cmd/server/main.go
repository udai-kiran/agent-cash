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

	"github.com/udai-kiran/agentic-cash/internal/application/service"
	"github.com/udai-kiran/agentic-cash/internal/config"
	"github.com/udai-kiran/agentic-cash/internal/infrastructure/auth"
	"github.com/udai-kiran/agentic-cash/internal/infrastructure/persistence/postgres"
	httpRouter "github.com/udai-kiran/agentic-cash/internal/interfaces/http"
	"github.com/udai-kiran/agentic-cash/internal/interfaces/http/handler"
)

func main() {
	// Load configuration
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create context for database initialization
	ctx := context.Background()

	// Initialize database connection pool
	pool, err := postgres.NewPool(ctx, &cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	log.Println("Connected to PostgreSQL successfully")

	// Initialize application tables
	if err := postgres.InitializeAppTables(ctx, pool); err != nil {
		log.Fatalf("Failed to initialize app tables: %v", err)
	}

	log.Println("Application tables initialized")

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
		log.Printf("Starting server on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
