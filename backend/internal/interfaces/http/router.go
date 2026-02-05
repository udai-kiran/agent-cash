package http

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/udai-kiran/agentic-cash/internal/infrastructure/auth"
	"github.com/udai-kiran/agentic-cash/internal/interfaces/http/handler"
	"github.com/udai-kiran/agentic-cash/internal/interfaces/http/middleware"
)

// RouterConfig holds dependencies for router setup
type RouterConfig struct {
	AccountHandler     *handler.AccountHandler
	AuthHandler        *handler.AuthHandler
	TransactionHandler *handler.TransactionHandler
	AnalyticsHandler   *handler.AnalyticsHandler
	CommodityHandler   *handler.CommodityHandler
	JWTManager         *auth.JWTManager
	AllowedOrigins     []string
}

// Router sets up the HTTP router
func Router(cfg *RouterConfig) *gin.Engine {
	r := gin.Default()

	// Apply middleware
	r.Use(middleware.RequestID())
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.RequestSizeLimit(10 * 1024 * 1024)) // 10MB limit
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.CORS(&middleware.CORSConfig{
		AllowedOrigins: cfg.AllowedOrigins,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Rate limiter for auth endpoints: 5 requests per minute
	authRateLimiter := middleware.NewRateLimiter(5, time.Minute)

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Auth routes (public, rate limited)
		auth := v1.Group("/auth")
		auth.Use(authRateLimiter.Middleware())
		{
			auth.POST("/register", cfg.AuthHandler.Register)
			auth.POST("/login", cfg.AuthHandler.Login)
			auth.POST("/refresh", cfg.AuthHandler.RefreshToken)
			auth.POST("/logout", cfg.AuthHandler.Logout)
		}

		// Account routes (public for demo, can be protected with middleware)
		accounts := v1.Group("/accounts")
		{
			accounts.GET("", cfg.AccountHandler.GetAccounts)
			accounts.GET("/hierarchy", cfg.AccountHandler.GetAccountHierarchy)
			accounts.GET("/:guid", cfg.AccountHandler.GetAccount)
			accounts.GET("/:guid/balance", cfg.AccountHandler.GetAccountBalance)
		}

		// Transaction routes (public for demo, can be protected with middleware)
		transactions := v1.Group("/transactions")
		{
			transactions.GET("", cfg.TransactionHandler.GetTransactions)
			transactions.GET("/:guid", cfg.TransactionHandler.GetTransaction)
		}

		// Commodity routes (public for demo, can be protected with middleware)
		commodities := v1.Group("/commodities")
		{
			commodities.GET("/currencies", cfg.CommodityHandler.GetCurrencies)
		}

		// Analytics routes (public for demo, can be protected with middleware)
		analytics := v1.Group("/analytics")
		{
			analytics.GET("/income-expense", cfg.AnalyticsHandler.GetIncomeExpense)
			analytics.GET("/category-breakdown", cfg.AnalyticsHandler.GetCategoryBreakdown)
			analytics.GET("/net-worth", cfg.AnalyticsHandler.GetNetWorth)
		}

		// Protected routes example
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTManager))
		{
			// Add protected routes here in the future
			// protected.GET("/profile", profileHandler.GetProfile)
		}
	}

	return r
}
