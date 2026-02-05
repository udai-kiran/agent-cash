package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/udai-kiran/agentic-cash/internal/application/dto"
)

// RequestSizeLimit returns a middleware that limits request body size
func RequestSizeLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)

		// Check if body exceeds limit
		if c.Request.ContentLength > maxBytes {
			c.JSON(http.StatusRequestEntityTooLarge, dto.ErrorResponse{
				Error:   "Request Entity Too Large",
				Message: "Request body exceeds maximum allowed size",
				Code:    http.StatusRequestEntityTooLarge,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
