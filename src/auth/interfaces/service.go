package authinterfaces

import (
	"context"

	"gobackend/src/auth/dto"
)

// AuthService encapsulates authentication business logic.
type AuthService interface {
	GetGoogleLoginURL(state string) string
	HandleGoogleCallback(ctx context.Context, req dto.GoogleCallbackRequest) (*dto.AuthResponse, error)
}
