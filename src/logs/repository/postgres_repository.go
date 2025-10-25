package repository

import (
	"context"
	"database/sql"

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

// EnsureSchema ensures the user_logs table exists.
func (r *PostgresRepository) EnsureSchema(ctx context.Context) error {
	const query = `
CREATE TABLE IF NOT EXISTS user_logs (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    action TEXT NOT NULL,
    detail TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('Asia/Jakarta', NOW())
);

CREATE INDEX IF NOT EXISTS user_logs_user_id_idx ON user_logs (user_id);
CREATE INDEX IF NOT EXISTS user_logs_created_at_idx ON user_logs (created_at DESC);
`
	_, err := r.db.ExecContext(ctx, query)
	return err
}

// FindAll retrieves logs using pagination parameters and returns total count.
func (r *PostgresRepository) FindAll(ctx context.Context, params pagination.Params) ([]dao.Log, int64, error) {
	const listQuery = `
SELECT l.id,
       l.user_id,
       COALESCE(u.name, ''),
       l.action,
       COALESCE(l.detail, ''),
       l.created_at
FROM user_logs l
LEFT JOIN users u ON u.id = l.user_id
ORDER BY l.created_at DESC
LIMIT $1 OFFSET $2
`
	rows, err := r.db.QueryContext(ctx, listQuery, params.Limit(), params.Offset())
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

	const countQuery = `SELECT COUNT(*) FROM user_logs`
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
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
