package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/udai-kiran/agentic-cash/pkg/logger"
)

// LoggingMiddleware logs HTTP requests with structured logging
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// Log after request completes
		duration := time.Since(start)
		statusCode := c.Writer.Status()
		requestID := GetRequestID(c)

		logFields := []any{
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status", statusCode),
			slog.Duration("duration", duration),
			slog.String("request_id", requestID),
			slog.String("client_ip", c.ClientIP()),
		}

		if len(c.Errors) > 0 {
			logFields = append(logFields, slog.String("error", c.Errors.String()))
			logger.Error("HTTP request completed with errors", logFields...)
		} else if statusCode >= 500 {
			logger.Error("HTTP request failed", logFields...)
		} else if statusCode >= 400 {
			logger.Warn("HTTP request client error", logFields...)
		} else {
			logger.Info("HTTP request completed", logFields...)
		}
	}
}
