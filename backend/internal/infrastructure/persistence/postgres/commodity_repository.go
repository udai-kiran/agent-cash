package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/udai-kiran/agentic-cash/internal/domain/entity"
	"github.com/udai-kiran/agentic-cash/internal/domain/repository"
)

// CommodityRepository implements repository.CommodityRepository for PostgreSQL
type CommodityRepository struct {
	db *pgxpool.Pool
}

// NewCommodityRepository creates a new PostgreSQL commodity repository
func NewCommodityRepository(db *pgxpool.Pool) repository.CommodityRepository {
	return &CommodityRepository{db: db}
}

// FindCurrencies retrieves all commodities with namespace 'CURRENCY'
func (r *CommodityRepository) FindCurrencies(ctx context.Context) ([]*entity.Commodity, error) {
	query := `SELECT guid, namespace, mnemonic, fullname, fraction
	          FROM commodities
	          WHERE namespace = 'CURRENCY'
	          ORDER BY mnemonic`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query currencies: %w", err)
	}
	defer rows.Close()

	var commodities []*entity.Commodity
	for rows.Next() {
		c := &entity.Commodity{}
		err := rows.Scan(&c.GUID, &c.Namespace, &c.Mnemonic, &c.Fullname, &c.Fraction)
		if err != nil {
			return nil, fmt.Errorf("failed to scan commodity: %w", err)
		}
		commodities = append(commodities, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating commodities: %w", err)
	}

	return commodities, nil
}

// FindByGUID retrieves a commodity by its GUID
func (r *CommodityRepository) FindByGUID(ctx context.Context, guid string) (*entity.Commodity, error) {
	query := `SELECT guid, namespace, mnemonic, fullname, fraction
	          FROM commodities
	          WHERE guid = $1`

	c := &entity.Commodity{}
	err := r.db.QueryRow(ctx, query, guid).Scan(&c.GUID, &c.Namespace, &c.Mnemonic, &c.Fullname, &c.Fraction)
	if err != nil {
		return nil, fmt.Errorf("failed to find commodity: %w", err)
	}

	return c, nil
}
