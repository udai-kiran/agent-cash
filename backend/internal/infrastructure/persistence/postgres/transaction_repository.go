package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/udai-kiran/agentic-cash/internal/domain/entity"
	"github.com/udai-kiran/agentic-cash/internal/domain/repository"
)

// TransactionRepository implements repository.TransactionRepository for PostgreSQL
type TransactionRepository struct {
	db *pgxpool.Pool
}

// NewTransactionRepository creates a new PostgreSQL transaction repository
func NewTransactionRepository(db *pgxpool.Pool) repository.TransactionRepository {
	return &TransactionRepository{db: db}
}

// FindAll retrieves all transactions with optional filtering
func (r *TransactionRepository) FindAll(ctx context.Context, filter *repository.TransactionFilter) ([]*entity.Transaction, error) {
	query := `
		SELECT DISTINCT t.guid, t.currency_guid, COALESCE(c.mnemonic, ''), t.num, t.post_date, t.enter_date, t.description
		FROM transactions t
		LEFT JOIN commodities c ON t.currency_guid = c.guid
	`

	var conditions []string
	var args []interface{}
	argPos := 1

	if filter != nil {
		if filter.AccountGUID != nil {
			conditions = append(conditions, fmt.Sprintf(`EXISTS (
				SELECT 1 FROM splits s WHERE s.tx_guid = t.guid AND s.account_guid = $%d
			)`, argPos))
			args = append(args, *filter.AccountGUID)
			argPos++
		}

		if filter.StartDate != nil {
			conditions = append(conditions, fmt.Sprintf("t.post_date >= $%d", argPos))
			args = append(args, *filter.StartDate)
			argPos++
		}

		if filter.EndDate != nil {
			conditions = append(conditions, fmt.Sprintf("t.post_date <= $%d", argPos))
			args = append(args, *filter.EndDate)
			argPos++
		}

		if filter.Description != nil {
			conditions = append(conditions, fmt.Sprintf("t.description ILIKE $%d", argPos))
			args = append(args, "%"+*filter.Description+"%")
			argPos++
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY t.post_date DESC, t.enter_date DESC"

	if filter != nil {
		if filter.Limit > 0 {
			query += fmt.Sprintf(" LIMIT $%d", argPos)
			args = append(args, filter.Limit)
			argPos++
		}

		if filter.Offset > 0 {
			query += fmt.Sprintf(" OFFSET $%d", argPos)
			args = append(args, filter.Offset)
		}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*entity.Transaction
	for rows.Next() {
		tx := &entity.Transaction{}
		err := rows.Scan(
			&tx.GUID,
			&tx.CurrencyGUID,
			&tx.CurrencyMnemonic,
			&tx.Num,
			&tx.PostDate,
			&tx.EnterDate,
			&tx.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		// Load splits for this transaction
		splits, err := r.loadSplitsForTransaction(ctx, tx.GUID)
		if err != nil {
			return nil, err
		}
		tx.Splits = splits

		transactions = append(transactions, tx)
	}

	return transactions, nil
}

// FindByGUID retrieves a transaction by its GUID
func (r *TransactionRepository) FindByGUID(ctx context.Context, guid string) (*entity.Transaction, error) {
	query := `
		SELECT t.guid, t.currency_guid, COALESCE(c.mnemonic, ''), t.num, t.post_date, t.enter_date, t.description
		FROM transactions t
		LEFT JOIN commodities c ON t.currency_guid = c.guid
		WHERE t.guid = $1
	`

	tx := &entity.Transaction{}
	err := r.db.QueryRow(ctx, query, guid).Scan(
		&tx.GUID,
		&tx.CurrencyGUID,
		&tx.CurrencyMnemonic,
		&tx.Num,
		&tx.PostDate,
		&tx.EnterDate,
		&tx.Description,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find transaction: %w", err)
	}

	// Load splits
	splits, err := r.loadSplitsForTransaction(ctx, guid)
	if err != nil {
		return nil, err
	}
	tx.Splits = splits

	return tx, nil
}

// FindByAccount retrieves transactions for a specific account
func (r *TransactionRepository) FindByAccount(ctx context.Context, accountGUID string, limit, offset int) ([]*entity.Transaction, error) {
	filter := &repository.TransactionFilter{
		AccountGUID: &accountGUID,
		Limit:       limit,
		Offset:      offset,
	}
	return r.FindAll(ctx, filter)
}

// Count returns the total number of transactions matching the filter
func (r *TransactionRepository) Count(ctx context.Context, filter *repository.TransactionFilter) (int64, error) {
	query := `SELECT COUNT(DISTINCT t.guid) FROM transactions t`

	var conditions []string
	var args []interface{}
	argPos := 1

	if filter != nil {
		if filter.AccountGUID != nil {
			conditions = append(conditions, fmt.Sprintf(`EXISTS (
				SELECT 1 FROM splits s WHERE s.tx_guid = t.guid AND s.account_guid = $%d
			)`, argPos))
			args = append(args, *filter.AccountGUID)
			argPos++
		}

		if filter.StartDate != nil {
			conditions = append(conditions, fmt.Sprintf("t.post_date >= $%d", argPos))
			args = append(args, *filter.StartDate)
			argPos++
		}

		if filter.EndDate != nil {
			conditions = append(conditions, fmt.Sprintf("t.post_date <= $%d", argPos))
			args = append(args, *filter.EndDate)
			argPos++
		}

		if filter.Description != nil {
			conditions = append(conditions, fmt.Sprintf("t.description ILIKE $%d", argPos))
			args = append(args, "%"+*filter.Description+"%")
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var count int64
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count transactions: %w", err)
	}

	return count, nil
}

// loadSplitsForTransaction loads splits for a transaction
func (r *TransactionRepository) loadSplitsForTransaction(ctx context.Context, txGUID string) ([]*entity.Split, error) {
	query := `
		SELECT s.guid, s.tx_guid, s.account_guid, s.memo, s.action,
		       s.reconcile_state, s.value_num, s.value_denom,
		       s.quantity_num, s.quantity_denom,
		       a.name as account_name, a.account_type
		FROM splits s
		LEFT JOIN accounts a ON s.account_guid = a.guid
		WHERE s.tx_guid = $1
		ORDER BY s.value_num DESC
	`

	rows, err := r.db.Query(ctx, query, txGUID)
	if err != nil {
		return nil, fmt.Errorf("failed to query splits: %w", err)
	}
	defer rows.Close()

	var splits []*entity.Split
	for rows.Next() {
		split := &entity.Split{
			Account: &entity.Account{},
		}
		err := rows.Scan(
			&split.GUID,
			&split.TxGUID,
			&split.AccountGUID,
			&split.Memo,
			&split.Action,
			&split.ReconcileState,
			&split.ValueNum,
			&split.ValueDenom,
			&split.QuantityNum,
			&split.QuantityDenom,
			&split.Account.Name,
			&split.Account.AccountType,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan split: %w", err)
		}
		split.Account.GUID = split.AccountGUID

		splits = append(splits, split)
	}

	return splits, nil
}
