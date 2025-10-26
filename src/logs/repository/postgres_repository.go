package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"gobackend/shared/pagination"
	"gobackend/src/logs/dao"
	loginterfaces "gobackend/src/logs/interfaces"
)

var _ loginterfaces.Repository = (*PostgresRepository)(nil)

// PostgresRepository implements user log queries against Postgres.
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new log repository.
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// EnsureSchema verifies that required tables and indexes exist.
func (r *PostgresRepository) EnsureSchema(ctx context.Context) error {
	const logsTableQuery = `
SELECT 1
FROM information_schema.tables
WHERE table_schema = 'public' AND table_name = 'user_logs'
`

	var exists int
	if err := r.db.QueryRowContext(ctx, logsTableQuery).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user_logs table not found; please run database migrations")
		}
		return err
	}

	const userIdxQuery = `
SELECT 1
FROM pg_indexes
WHERE schemaname = 'public' AND indexname = 'user_logs_user_id_idx'
`

	if err := r.db.QueryRowContext(ctx, userIdxQuery).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("index user_logs_user_id_idx not found; please run database migrations")
		}
		return err
	}

	const createdIdxQuery = `
SELECT 1
FROM pg_indexes
WHERE schemaname = 'public' AND indexname = 'user_logs_created_at_idx'
`

	if err := r.db.QueryRowContext(ctx, createdIdxQuery).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("index user_logs_created_at_idx not found; please run database migrations")
		}
		return err
	}

	return nil
}

// FindAll retrieves logs using pagination parameters and optional filters, returning the total count.
func (r *PostgresRepository) FindAll(ctx context.Context, params pagination.Params, userID *int64) ([]dao.Log, int64, error) {
	baseQuery := `
SELECT l.id,
       l.user_id,
       COALESCE(u.name, ''),
       l.action,
       COALESCE(l.detail, ''),
       l.created_at
FROM user_logs l
LEFT JOIN users u ON u.id = l.user_id`

	countQuery := `SELECT COUNT(*) FROM user_logs`

	var (
		whereClause string
		args        []interface{}
		countArgs   []interface{}
	)

	if userID != nil {
		whereClause = " WHERE l.user_id = $1"
		countQuery += " WHERE user_id = $1"
		args = append(args, *userID)
		countArgs = append(countArgs, *userID)
	}

	limitPlaceholder := len(args) + 1
	offsetPlaceholder := len(args) + 2
	args = append(args, params.Limit(), params.Offset())

	query := fmt.Sprintf("%s%s ORDER BY l.created_at DESC LIMIT $%d OFFSET $%d", baseQuery, whereClause, limitPlaceholder, offsetPlaceholder)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []dao.Log
	for rows.Next() {
		var (
			log      dao.Log
			userName string
		)

		if err := rows.Scan(&log.ID, &log.UserID, &userName, &log.Action, &log.Detail, &log.CreatedAt); err != nil {
			return nil, 0, err
		}
		log.UserName = userName
		logs = append(logs, log)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	var total int64
	if len(countArgs) > 0 {
		if err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total); err != nil {
			return nil, 0, err
		}
	} else {
		if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
			return nil, 0, err
		}
	}

	return logs, total, nil
}

// Create inserts a new log entry.
func (r *PostgresRepository) Create(ctx context.Context, entry dao.Log) error {
	const query = `
INSERT INTO user_logs (user_id, action, detail)
VALUES ($1, $2, $3)
`
	_, err := r.db.ExecContext(ctx, query, entry.UserID, entry.Action, entry.Detail)
	return err
}
