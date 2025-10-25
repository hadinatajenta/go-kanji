package interfaces

import (
	"context"

	"gobackend/shared/pagination"
	"gobackend/src/users/logs/dao"
)

// Repository describes persistence layer for user logs.
type Repository interface {
	FindAll(ctx context.Context, params pagination.Params) ([]dao.Log, int64, error)
	EnsureSchema(ctx context.Context) error
}
