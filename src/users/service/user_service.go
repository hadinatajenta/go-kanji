package service

import (
	"context"

	"gobackend/src/users/dto"
	userinterfaces "gobackend/src/users/interfaces"
)

var _ userinterfaces.UserService = (*UserServiceImpl)(nil)

// UserServiceImpl provides user read operations.
type UserServiceImpl struct {
	repo userinterfaces.UserRepository
}

// NewUserService creates a new UserServiceImpl instance.
func NewUserService(repo userinterfaces.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{repo: repo}
}

// ListUsers retrieves all users and maps them into DTOs.
func (s *UserServiceImpl) ListUsers(ctx context.Context) ([]dto.User, error) {
	users, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]dto.User, 0, len(users))
	for _, user := range users {
		result = append(result, dto.User{
			ID:          user.ID,
			Email:       user.Email,
			Name:        user.Name,
			PictureURL:  user.PictureURL,
			LastLoginAt: user.LastLoginAt,
			CreatedAt:   user.CreatedAt,
			Provider:    user.Provider,
		})
	}

	return result, nil
}
