package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/udai-kiran/agentic-cash/internal/domain/entity"
	"github.com/udai-kiran/agentic-cash/internal/domain/repository"
	"github.com/udai-kiran/agentic-cash/pkg/gnucash"
)

// AccountRepository implements repository.AccountRepository for PostgreSQL
type AccountRepository struct {
	db *pgxpool.Pool
}

// NewAccountRepository creates a new PostgreSQL account repository
func NewAccountRepository(db *pgxpool.Pool) repository.AccountRepository {
	return &AccountRepository{db: db}
}

const accountSelectColumns = `a.guid, a.name, a.account_type, a.commodity_guid, a.commodity_scu,
		       a.parent_guid, a.code, a.description, a.hidden, a.placeholder,
		       COALESCE(c.mnemonic, '')`

// scanAccount scans a row into an Account entity, handling int-to-bool conversion
// for hidden and placeholder columns (GnuCash stores these as integer 0/1).
func scanAccount(row pgx.Row) (*entity.Account, error) {
	account := &entity.Account{}
	var hidden, placeholder int
	err := row.Scan(
		&account.GUID,
		&account.Name,
		&account.AccountType,
		&account.CommodityGUID,
		&account.CommoditySCU,
		&account.ParentGUID,
		&account.Code,
		&account.Description,
		&hidden,
		&placeholder,
		&account.CommodityMnemonic,
	)
	if err != nil {
		return nil, err
	}
	account.Hidden = hidden != 0
	account.Placeholder = placeholder != 0
	return account, nil
}

// FindAll retrieves all accounts
func (r *AccountRepository) FindAll(ctx context.Context) ([]*entity.Account, error) {
	query := fmt.Sprintf(`SELECT %s FROM accounts a LEFT JOIN commodities c ON a.commodity_guid = c.guid ORDER BY a.name`, accountSelectColumns)

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query accounts: %w", err)
	}
	defer rows.Close()

	var accounts []*entity.Account
	for rows.Next() {
		account, err := scanAccount(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan account: %w", err)
		}
		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating accounts: %w", err)
	}

	return accounts, nil
}

// FindByGUID retrieves an account by its GUID
func (r *AccountRepository) FindByGUID(ctx context.Context, guid string) (*entity.Account, error) {
	query := fmt.Sprintf(`SELECT %s FROM accounts a LEFT JOIN commodities c ON a.commodity_guid = c.guid WHERE a.guid = $1`, accountSelectColumns)

	account, err := scanAccount(r.db.QueryRow(ctx, query, guid))
	if err != nil {
		return nil, fmt.Errorf("failed to find account: %w", err)
	}

	return account, nil
}

// FindHierarchy retrieves the complete account hierarchy
func (r *AccountRepository) FindHierarchy(ctx context.Context) ([]*entity.Account, error) {
	accounts, err := r.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// Build a map for quick lookup
	accountMap := make(map[string]*entity.Account)
	for _, account := range accounts {
		accountMap[account.GUID] = account
		account.Children = []*entity.Account{}
	}

	// Build the hierarchy
	var roots []*entity.Account
	for _, account := range accounts {
		if account.ParentGUID == nil || *account.ParentGUID == "" {
			roots = append(roots, account)
		} else {
			parent, exists := accountMap[*account.ParentGUID]
			if exists {
				parent.Children = append(parent.Children, account)
			}
		}
	}

	return roots, nil
}

// FindByType retrieves accounts by type
func (r *AccountRepository) FindByType(ctx context.Context, accountType entity.AccountType) ([]*entity.Account, error) {
	query := fmt.Sprintf(`SELECT %s FROM accounts a LEFT JOIN commodities c ON a.commodity_guid = c.guid WHERE a.account_type = $1 ORDER BY a.name`, accountSelectColumns)

	rows, err := r.db.Query(ctx, query, accountType)
	if err != nil {
		return nil, fmt.Errorf("failed to query accounts by type: %w", err)
	}
	defer rows.Close()

	var accounts []*entity.Account
	for rows.Next() {
		account, err := scanAccount(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan account: %w", err)
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

// GetBalance calculates the current balance for an account
func (r *AccountRepository) GetBalance(ctx context.Context, guid string) (int64, int64, error) {
	query := `
		SELECT COALESCE(SUM(s.quantity_num), 0) as total_num,
		       COALESCE(MAX(s.quantity_denom), 100) as denom
		FROM splits s
		WHERE s.account_guid = $1
	`

	var numerator, denominator int64
	err := r.db.QueryRow(ctx, query, guid).Scan(&numerator, &denominator)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to calculate balance: %w", err)
	}

	return numerator, denominator, nil
}

// GetBalanceWithChildren calculates the balance including child accounts
func (r *AccountRepository) GetBalanceWithChildren(ctx context.Context, guid string) (int64, int64, error) {
	// First get the account to check if it's a debit or credit account
	account, err := r.FindByGUID(ctx, guid)
	if err != nil {
		return 0, 0, err
	}

	// Recursive CTE to get all child accounts
	query := `
		WITH RECURSIVE account_tree AS (
			SELECT guid FROM accounts WHERE guid = $1
			UNION ALL
			SELECT a.guid FROM accounts a
			INNER JOIN account_tree at ON a.parent_guid = at.guid
		)
		SELECT COALESCE(SUM(s.quantity_num), 0) as total_num,
		       COALESCE(MAX(s.quantity_denom), 100) as denom
		FROM splits s
		WHERE s.account_guid IN (SELECT guid FROM account_tree)
	`

	var numerator, denominator int64
	err = r.db.QueryRow(ctx, query, guid).Scan(&numerator, &denominator)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to calculate balance with children: %w", err)
	}

	// Normalize the sign based on account type
	numerator = gnucash.NormalizeSign(numerator, account.IsDebitAccount())

	return numerator, denominator, nil
}
