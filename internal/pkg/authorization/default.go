package authorization

import (
	"errors"
	"strings"

	"github.com/agrim123/gatekeeper/internal/pkg/store"
)

type DefaultModule struct{}

func NewDefaultModule() *DefaultModule {
	return &DefaultModule{}
}

func (dm DefaultModule) IsAuthorized(user, plan, option string) (bool, error) {
	userRoles := make(map[string]bool)

	allowedPlans := make(map[string]bool)
	for _, role := range store.Users[user].Roles {
		for _, p := range store.Roles[role].AllowedPlans {
			allowedPlans[p] = true
		}
	}

	for p := range allowedPlans {
		if p == "*" {
			return true, nil
		}

		// TODO: is regex better?
		if strings.Contains(p, ".*") {
			for opt := range store.Plans[plan].Opts {
				userRoles[plan+"."+opt] = true
			}

			continue
		}

		userRoles[p] = true
	}

	rolesToUse := []string{
		plan + "." + option,
	}

	for _, role := range rolesToUse {
		if _, ok := userRoles[role]; !ok {
			return false, errors.New("Invalid role")
		}
	}

	return true, nil
}
