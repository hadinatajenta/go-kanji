package interfaces

import (
	"context"

	"gobackend/src/users/dto"
)

// UserService exposes user-related business logic.
type UserService interface {
	ListUsers(ctx context.Context) ([]dto.User, error)
}
