package repository

import (
	"context"

	"github.com/udai-kiran/agentic-cash/internal/domain/entity"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *entity.User) error

	// FindByEmail retrieves a user by email
	FindByEmail(ctx context.Context, email string) (*entity.User, error)

	// FindByID retrieves a user by ID
	FindByID(ctx context.Context, id int64) (*entity.User, error)

	// Update updates a user
	Update(ctx context.Context, user *entity.User) error

	// CreateRefreshToken stores a refresh token
	CreateRefreshToken(ctx context.Context, userID int64, token string, expiresAt int64) error

	// ValidateRefreshToken checks if a refresh token is valid and returns the user ID
	ValidateRefreshToken(ctx context.Context, token string) (int64, error)

	// DeleteRefreshToken removes a refresh token
	DeleteRefreshToken(ctx context.Context, token string) error

	// DeleteUserRefreshTokens removes all refresh tokens for a user
	DeleteUserRefreshTokens(ctx context.Context, userID int64) error
}
