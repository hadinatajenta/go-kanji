package service

import (
	"context"
	"strings"

	"gobackend/shared/identity"
	"gobackend/src/users/dto"
	userinterfaces "gobackend/src/users/interfaces"
)

var _ userinterfaces.UserService = (*UserServiceImpl)(nil)

// UserServiceImpl provides user read operations.
type UserServiceImpl struct {
	repo       userinterfaces.UserRepository
	refEncoder *identity.UserReferenceEncoder
}

// NewUserService creates a new UserServiceImpl instance.
func NewUserService(repo userinterfaces.UserRepository, refEncoder *identity.UserReferenceEncoder) *UserServiceImpl {
	return &UserServiceImpl{repo: repo, refEncoder: refEncoder}
}

// ListUsers retrieves all users and maps them into DTOs.
func (s *UserServiceImpl) ListUsers(ctx context.Context) ([]dto.User, error) {
	users, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]dto.User, 0, len(users))
	for _, user := range users {
		reference, refErr := s.refEncoder.Encode(user.ID)
		if refErr != nil {
			return nil, refErr
		}

		maskedEmail := maskEmail(user.Email)

		result = append(result, dto.User{
			Reference:   reference,
			Email:       maskedEmail,
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
