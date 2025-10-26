package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gobackend/src/auth/dao"
	authinterfaces "gobackend/src/auth/interfaces"
)

var _ authinterfaces.UserRepository = (*PostgresUserRepository)(nil)

// PostgresUserRepository persists user records in Postgres.
type PostgresUserRepository struct {
	db          *sql.DB
	nowProvider func() time.Time
}

// NewPostgresUserRepository constructs a PostgresUserRepository and ensures the expected schema exists.
func NewPostgresUserRepository(db *sql.DB) (*PostgresUserRepository, error) {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return nil, err
	}

	repo := &PostgresUserRepository{
		db: db,
		nowProvider: func() time.Time {
			return time.Now().In(location)
		},
	}
	if err := repo.ensureSchema(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *PostgresUserRepository) ensureSchema() error {
	const usersTableQuery = `
SELECT 1
FROM information_schema.tables
WHERE table_schema = 'public' AND table_name = 'users'
`

	var exists int
	if err := r.db.QueryRow(usersTableQuery).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("users table not found; please run database migrations: %w", err)
		}
		return err
	}

	const usersIndexQuery = `
SELECT 1
FROM pg_indexes
WHERE schemaname = 'public' AND indexname = 'users_provider_provider_id_idx'
`

	if err := r.db.QueryRow(usersIndexQuery).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("index users_provider_provider_id_idx not found; please run database migrations: %w", err)
		}
		return err
	}

	return nil
}

// FindByProvider locates a user by provider details.
func (r *PostgresUserRepository) FindByProvider(ctx context.Context, provider, providerID string) (*dao.User, error) {
	const query = `
SELECT id, email, name, provider, provider_id, picture_url, created_at, last_login_at
FROM users
WHERE provider = $1 AND provider_id = $2
LIMIT 1
`

	row := r.db.QueryRowContext(ctx, query, provider, providerID)

	var (
		user      dao.User
		lastLogin sql.NullTime
	)

	if err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Provider,
		&user.ProviderID,
		&user.PictureURL,
		&user.CreatedAt,
		&lastLogin,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if lastLogin.Valid {
		user.LastLoginAt = lastLogin.Time
	}

	return &user, nil
}

// Create inserts a new user record.
func (r *PostgresUserRepository) Create(ctx context.Context, user dao.User) (*dao.User, error) {
	const query = `
INSERT INTO users (email, name, provider, provider_id, picture_url, created_at, last_login_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, email, name, provider, provider_id, picture_url, created_at, last_login_at
`

	now := r.nowProvider()
	lastLogin := now

	var result dao.User
	var lastLoginTime sql.NullTime

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.Name,
		user.Provider,
		user.ProviderID,
		user.PictureURL,
		now,
		lastLogin,
	).Scan(
		&result.ID,
		&result.Email,
		&result.Name,
		&result.Provider,
		&result.ProviderID,
		&result.PictureURL,
		&result.CreatedAt,
		&lastLoginTime,
	)
	if err != nil {
		return nil, err
	}

	if lastLoginTime.Valid {
		result.LastLoginAt = lastLoginTime.Time
	} else {
		result.LastLoginAt = now
	}

	return &result, nil
}

// UpdateLoginTimestamp refreshes the user's last_login_at column.
func (r *PostgresUserRepository) UpdateLoginTimestamp(ctx context.Context, userID int64) error {
	const query = `
UPDATE users
SET last_login_at = $1
WHERE id = $2
`

	_, err := r.db.ExecContext(ctx, query, r.nowProvider(), userID)
	return err
}
