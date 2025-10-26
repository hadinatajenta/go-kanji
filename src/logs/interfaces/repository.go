package interfaces

import (
	"context"

	"gobackend/shared/pagination"
	"gobackend/src/logs/dao"
)

// Repository describes persistence layer for user logs.
type Repository interface {
	FindAll(ctx context.Context, params pagination.Params, userID *int64) ([]dao.Log, int64, error)
	Create(ctx context.Context, entry dao.Log) error
	EnsureSchema(ctx context.Context) error
}
