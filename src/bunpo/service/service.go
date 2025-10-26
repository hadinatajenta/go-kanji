package service

import (
	"context"

	bunpointerfaces "gobackend/src/bunpo/interfaces"
)

var _ bunpointerfaces.Service = (*bunpoService)(nil)

type bunpoService struct{}

// NewService constructs a bunpo service implementation.
func NewService() bunpointerfaces.Service {
	return &bunpoService{}
}

// Test returns a simple success message.
func (s *bunpoService) Test(ctx context.Context) (string, error) {
	return "endpoint success", nil
}
