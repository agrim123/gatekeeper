package authorization

import "context"

type Module interface {
	IsAuthorized(ctx context.Context, plan, option string) (bool, error)
}
