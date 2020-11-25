package authentication

import "context"

type Module interface {
	IsAuthenticated(ctx context.Context) (bool, error)
}
