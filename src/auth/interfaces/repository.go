package authinterfaces

import (
	"context"

	"gobackend/src/auth/dao"
)

// UserRepository describes storage operations required by the auth service.
type UserRepository interface {
	FindByProvider(ctx context.Context, provider, providerID string) (*dao.User, error)
	Create(ctx context.Context, user dao.User) (*dao.User, error)
	UpdateLoginTimestamp(ctx context.Context, userID int64) error
}
