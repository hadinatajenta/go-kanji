package bunpointerfaces

import "context"

// Service describes bunpo feature business logic.
type Service interface {
	Test(ctx context.Context) (string, error)
}
