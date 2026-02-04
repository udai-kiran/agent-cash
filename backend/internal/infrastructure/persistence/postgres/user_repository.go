package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/udai-kiran/agentic-cash/internal/domain/entity"
	"github.com/udai-kiran/agentic-cash/internal/domain/repository"
)

// UserRepository implements repository.UserRepository for PostgreSQL
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new PostgreSQL user repository
func NewUserRepository(db *pgxpool.Pool) repository.UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO app_users (email, password_hash, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, user.Email, user.PasswordHash).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// FindByEmail retrieves a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, created_at, updated_at
		FROM app_users
		WHERE email = $1
	`

	user := &entity.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

// FindByID retrieves a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, created_at, updated_at
		FROM app_users
		WHERE id = $1
	`

	user := &entity.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE app_users
		SET email = $1, password_hash = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query, user.Email, user.PasswordHash, user.ID).Scan(&user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// CreateRefreshToken stores a refresh token
func (r *UserRepository) CreateRefreshToken(ctx context.Context, userID int64, token string, expiresAt int64) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, to_timestamp($3))
	`

	_, err := r.db.Exec(ctx, query, userID, token, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to create refresh token: %w", err)
	}

	return nil
}

// ValidateRefreshToken checks if a refresh token is valid
func (r *UserRepository) ValidateRefreshToken(ctx context.Context, token string) (int64, error) {
	query := `
		SELECT user_id FROM refresh_tokens
		WHERE token = $1 AND expires_at > NOW()
	`

	var userID int64
	err := r.db.QueryRow(ctx, query, token).Scan(&userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, fmt.Errorf("invalid or expired token")
		}
		return 0, fmt.Errorf("failed to validate token: %w", err)
	}

	return userID, nil
}

// DeleteRefreshToken removes a refresh token
func (r *UserRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = $1`

	_, err := r.db.Exec(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	return nil
}

// DeleteUserRefreshTokens removes all refresh tokens for a user
func (r *UserRepository) DeleteUserRefreshTokens(ctx context.Context, userID int64) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`

	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user refresh tokens: %w", err)
	}

	return nil
}
