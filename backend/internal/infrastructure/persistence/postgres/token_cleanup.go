package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/udai-kiran/agentic-cash/pkg/logger"
)

// TokenCleanupService handles periodic cleanup of expired refresh tokens
type TokenCleanupService struct {
	db       *pgxpool.Pool
	interval time.Duration
	stopChan chan struct{}
}

// NewTokenCleanupService creates a new token cleanup service
func NewTokenCleanupService(db *pgxpool.Pool, interval time.Duration) *TokenCleanupService {
	return &TokenCleanupService{
		db:       db,
		interval: interval,
		stopChan: make(chan struct{}),
	}
}

// Start begins the periodic cleanup process
func (s *TokenCleanupService) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	logger.Info("Token cleanup service started", "interval", s.interval)

	// Run cleanup immediately on start
	s.cleanup(ctx)

	for {
		select {
		case <-ticker.C:
			s.cleanup(ctx)
		case <-s.stopChan:
			logger.Info("Token cleanup service stopped")
			return
		case <-ctx.Done():
			logger.Info("Token cleanup service context cancelled")
			return
		}
	}
}

// Stop halts the cleanup service
func (s *TokenCleanupService) Stop() {
	close(s.stopChan)
}

// cleanup removes expired tokens from the database
func (s *TokenCleanupService) cleanup(ctx context.Context) {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`

	result, err := s.db.Exec(ctx, query)
	if err != nil {
		logger.Error("Failed to cleanup expired tokens", "error", err)
		return
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected > 0 {
		logger.Info("Cleaned up expired refresh tokens", "count", rowsAffected)
	}
}
