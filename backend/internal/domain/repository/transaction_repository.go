package repository

import (
	"context"
	"time"

	"github.com/udai-kiran/agentic-cash/internal/domain/entity"
)

// TransactionFilter defines filtering criteria for transactions
type TransactionFilter struct {
	AccountGUID  *string
	StartDate    *time.Time
	EndDate      *time.Time
	Description  *string
	MinAmount    *int64
	MaxAmount    *int64
	Limit        int
	Offset       int
}

// AccountAggregate represents aggregated data for an account
type AccountAggregate struct {
	AccountGUID string
	AccountName string
	TotalAmount int64  // Numerator in rational representation
	Denominator int64  // Denominator in rational representation
	Count       int
}

// TransactionRepository defines the interface for transaction data access
type TransactionRepository interface {
	// FindAll retrieves all transactions with optional filtering
	FindAll(ctx context.Context, filter *TransactionFilter) ([]*entity.Transaction, error)

	// FindByGUID retrieves a transaction by its GUID
	FindByGUID(ctx context.Context, guid string) (*entity.Transaction, error)

	// FindByAccount retrieves transactions for a specific account
	FindByAccount(ctx context.Context, accountGUID string, limit, offset int) ([]*entity.Transaction, error)

	// Count returns the total number of transactions matching the filter
	Count(ctx context.Context, filter *TransactionFilter) (int64, error)

	// AggregateByAccountType returns aggregated transaction data grouped by account for accounts of specified type
	AggregateByAccountType(ctx context.Context, accountType entity.AccountType, startDate, endDate *time.Time) ([]*AccountAggregate, error)
}
