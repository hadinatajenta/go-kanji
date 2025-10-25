package service

import (
	"context"
	"strings"

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
			Email:       maskEmail(user.Email),
			Name:        user.Name,
			PictureURL:  user.PictureURL,
			LastLoginAt: user.LastLoginAt,
			CreatedAt:   user.CreatedAt,
			Provider:    user.Provider,
		})
	}

	return result, nil
}

func maskEmail(email string) string {
	const maskedSegment = "*****"

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	local, domain := parts[0], parts[1]
	prefix := local
	if len(prefix) > 3 {
		prefix = prefix[:3]
	}

	return prefix + maskedSegment + "@" + domain
}
