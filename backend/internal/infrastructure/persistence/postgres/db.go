package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/udai-kiran/agentic-cash/internal/config"
)

// NewPool creates a new PostgreSQL connection pool
func NewPool(ctx context.Context, cfg *config.DatabaseConfig) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to parse pool config: %w", err)
	}

	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

// InitializeAppTables creates application-specific tables if they don't exist
func InitializeAppTables(ctx context.Context, pool *pgxpool.Pool) error {
	// Create users table
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS app_users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create app_users table: %w", err)
	}

	// Create refresh tokens table
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS refresh_tokens (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES app_users(id) ON DELETE CASCADE,
			token VARCHAR(500) UNIQUE NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create refresh_tokens table: %w", err)
	}

	// Create index on user_id for refresh_tokens
	_, err = pool.Exec(ctx, `
		CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id)
	`)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	return nil
}
