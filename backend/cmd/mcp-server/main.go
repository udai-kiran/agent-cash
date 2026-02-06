package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/udai-kiran/agentic-cash/internal/application/service"
	"github.com/udai-kiran/agentic-cash/internal/config"
	"github.com/udai-kiran/agentic-cash/internal/infrastructure/persistence/postgres"
	"github.com/udai-kiran/agentic-cash/internal/infrastructure/mcp"
	"github.com/udai-kiran/agentic-cash/pkg/logger"
)

func main() {
	// Initialize logger
	logger.Init(os.Getenv("GO_ENV") == "production")

	// Create context for database initialization
	ctx := context.Background()

	// Build database config from environment variables (MCP server doesn't use config file)
	dbConfig := &config.DatabaseConfig{
		Host:     getEnvOrDefault("DATABASE_HOST", "localhost"),
		Port:     getEnvAsInt("DATABASE_PORT", 5432),
		User:     getEnvOrDefault("DATABASE_USER", "gnucash"),
		Password: getEnvOrDefault("DATABASE_PASSWORD", "gnucash_password"),
		DBName:   getEnvOrDefault("DATABASE_NAME", "gnucash"),
		SSLMode:  getEnvOrDefault("DATABASE_SSLMODE", "disable"),
		MaxConns: 10,
		MinConns: 2,
	}

	// Initialize database connection pool
	pool, err := postgres.NewPool(ctx, dbConfig)
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

	// Initialize repositories
	accountRepo := postgres.NewAccountRepository(pool)
	transactionRepo := postgres.NewTransactionRepository(pool)
	commodityRepo := postgres.NewCommodityRepository(pool)

	// Initialize services
	analyticsService := service.NewAnalyticsService(accountRepo, transactionRepo)

	// Create MCP server
	mcpServer := mcp.NewMCPServer(
		accountRepo,
		transactionRepo,
		commodityRepo,
		analyticsService,
	)

	logger.Info("MCP server initialized successfully")

	// Set up signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start MCP server in a goroutine
	go func() {
		logger.Info("Starting MCP server...")
		if err := mcpServer.Start(ctx); err != nil {
			logger.Error("MCP server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	<-quit

	logger.Info("Shutting down MCP server...")

	// Graceful shutdown context
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the MCP server
	if err := mcpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("Error during MCP server shutdown", "error", err)
	}

	logger.Info("MCP server exited")
}

// Helper functions for environment variables
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		if _, err := fmt.Sscanf(value, "%d", &intValue); err == nil {
			return intValue
		}
	}
	return defaultValue
}
