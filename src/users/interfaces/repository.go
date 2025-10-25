package interfaces

import (
	"context"

	"gobackend/src/users/dao"
)

// UserRepository describes persistence operations used by the user service.
type UserRepository interface {
	FindAll(ctx context.Context) ([]dao.User, error)
}
