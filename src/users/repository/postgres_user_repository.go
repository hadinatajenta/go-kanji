package repository

import (
	"context"
	"database/sql"

	"gobackend/src/users/dao"
	userinterfaces "gobackend/src/users/interfaces"
)

var _ userinterfaces.UserRepository = (*PostgresUserRepository)(nil)

// PostgresUserRepository reads user records from Postgres.
type PostgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository builds a new PostgresUserRepository.
func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// FindAll returns all users ordered by creation date descending.
func (r *PostgresUserRepository) FindAll(ctx context.Context) ([]dao.User, error) {
	const query = `
SELECT id, email, name, provider, provider_id, picture_url, created_at, COALESCE(last_login_at, created_at)
FROM users
ORDER BY created_at DESC
`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []dao.User

	for rows.Next() {
		var user dao.User
		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.Provider,
			&user.ProviderID,
			&user.PictureURL,
			&user.CreatedAt,
			&user.LastLoginAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
