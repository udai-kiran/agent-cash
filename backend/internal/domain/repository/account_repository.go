package repository

import (
	"context"

	"github.com/udai-kiran/agentic-cash/internal/domain/entity"
)

// AccountRepository defines the interface for account data access
type AccountRepository interface {
	// FindAll retrieves all accounts
	FindAll(ctx context.Context) ([]*entity.Account, error)

	// FindByGUID retrieves an account by its GUID
	FindByGUID(ctx context.Context, guid string) (*entity.Account, error)

	// FindHierarchy retrieves the complete account hierarchy
	FindHierarchy(ctx context.Context) ([]*entity.Account, error)

	// FindByType retrieves accounts by type
	FindByType(ctx context.Context, accountType entity.AccountType) ([]*entity.Account, error)

	// GetBalance calculates the current balance for an account
	GetBalance(ctx context.Context, guid string) (int64, int64, error)
}
