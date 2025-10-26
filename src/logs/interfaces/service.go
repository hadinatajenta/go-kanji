package interfaces

import (
	"context"

	"gobackend/shared/pagination"
	"gobackend/src/logs/dto"
)

// Service defines operations for user logs.
type Service interface {
	ListLogs(ctx context.Context, params pagination.Params, userID *int64) ([]dto.Log, int64, error)
	Record(ctx context.Context, entry dto.NewLog) error
}
