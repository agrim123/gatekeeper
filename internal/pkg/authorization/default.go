package authorization

import (
	"context"
	"fmt"
	"strings"

	"github.com/agrim123/gatekeeper/internal/constants"
	"github.com/agrim123/gatekeeper/internal/pkg/store"
)

type DefaultModule struct{}

func NewDefaultModule() *DefaultModule {
	return &DefaultModule{}
}

func (dm DefaultModule) IsAuthorized(ctx context.Context, plan, option string) (bool, error) {
	if plan == "" {
		return false, fmt.Errorf("No plan specfied. Available plans: %v", store.GetAvailablePlans())
	}

	if _, ok := store.Plans[plan]; !ok {
		return false, fmt.Errorf("Invalid plan specfied. Available plans: %v", store.GetAvailablePlans())
	}

	allowedOptions := make([]string, 0)
	for _, role := range store.Users[ctx.Value(constants.UserContextKey).(string)].Roles {
		for _, p := range store.Roles[role].AllowedPlans {
			if p == "*" {
				return true, nil
			}

			if strings.Contains(p, ".*") {
				for opt := range store.Plans[plan].Opts {
					allowedOptions = append(allowedOptions, opt)
				}

				continue
			}

			allowedPlan := strings.Split(p, ".")
			if len(allowedPlan) != 2 {
				continue
			}

			if allowedPlan[0] == plan {
				allowedOptions = append(allowedOptions, allowedPlan[1])
			}
		}
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
