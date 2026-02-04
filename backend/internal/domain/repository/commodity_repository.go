package repository

import (
	"context"

	"github.com/udai-kiran/agentic-cash/internal/domain/entity"
)

// CommodityRepository defines the interface for commodity data access
type CommodityRepository interface {
	// FindCurrencies retrieves all commodities with namespace 'CURRENCY'
	FindCurrencies(ctx context.Context) ([]*entity.Commodity, error)

	// FindByGUID retrieves a commodity by its GUID
	FindByGUID(ctx context.Context, guid string) (*entity.Commodity, error)
}
