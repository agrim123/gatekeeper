package authorization

import (
	"context"
	"fmt"

	"github.com/agrim123/gatekeeper/internal/constants"
	"github.com/agrim123/gatekeeper/internal/store"
)

type DefaultModule struct {
	ctx context.Context
}

func NewDefaultModule(ctx context.Context) *DefaultModule {
	return &DefaultModule{
		ctx: ctx,
	}
}

func (dm DefaultModule) IsAuthorized(plan, option string) (bool, error) {
	if plan == "" {
		return false, fmt.Errorf("No plan specfied. Available plans: %v", store.Store.GetAvailablePlans())
	}

	if _, ok := store.Store.Plans[plan]; !ok {
		return false, fmt.Errorf("Invalid plan specfied. Available plans: %v", store.Store.GetAvailablePlans())
	}

	allowedOptions := make([]string, 0)
	if value, ok := store.Store.GetAllowedCommandsForUser(dm.ctx.Value(constants.UserContextKey).(string))[plan]; ok {
		allowedOptions = value
	}

	if option == "" {
		return false, fmt.Errorf("No option specified. Available options: %v", allowedOptions)
	}

	found := false
	for _, opt := range allowedOptions {
		if opt == option {
			found = true
			break
		}
	}

	if !found {
		return false, fmt.Errorf("Not authorized to run the specified plan: %s %s. Allowed options: %v", plan, option, allowedOptions)
	}

	return true, nil
}
