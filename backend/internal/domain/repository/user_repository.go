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
}
