package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/udai-kiran/agentic-cash/internal/application/dto"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*bucket
	rate     int           // tokens per window
	window   time.Duration // time window
	cleanupTicker *time.Ticker
}

type bucket struct {
	tokens    int
	lastReset time.Time
}

// NewRateLimiter creates a new rate limiter
// rate: number of requests allowed per window
// window: time window duration
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		buckets: make(map[string]*bucket),
		rate:    rate,
		window:  window,
		cleanupTicker: time.NewTicker(window * 2),
	}

	// Cleanup old buckets periodically
	go func() {
		for range rl.cleanupTicker.C {
			rl.cleanup()
		}
	}()

	return rl
}

// cleanup removes old buckets
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for key, b := range rl.buckets {
		if now.Sub(b.lastReset) > rl.window*3 {
			delete(rl.buckets, key)
		}
	}
}

// Allow checks if a request should be allowed
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	b, exists := rl.buckets[key]

	if !exists || now.Sub(b.lastReset) > rl.window {
		// Create new bucket or reset existing
		rl.buckets[key] = &bucket{
			tokens:    rl.rate - 1,
			lastReset: now,
		}
		return true
	}

	if b.tokens > 0 {
		b.tokens--
		return true
	}

	return false
}

// Middleware returns a Gin middleware that rate limits requests
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use IP address as the key
		key := c.ClientIP()

		if !rl.Allow(key) {
			c.JSON(http.StatusTooManyRequests, dto.ErrorResponse{
				Error:   "Rate Limit Exceeded",
				Message: "Too many requests. Please try again later.",
				Code:    http.StatusTooManyRequests,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Stop stops the cleanup ticker
func (rl *RateLimiter) Stop() {
	rl.cleanupTicker.Stop()
}
