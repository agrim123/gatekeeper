package authorization

import "context"

type Module interface {
	IsAuthorized(ctx context.Context) (bool, error)
}
