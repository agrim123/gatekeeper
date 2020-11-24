package authorization

import (
	"errors"
	"strings"

	"github.com/agrim123/gatekeeper/internal/pkg/root"
)

func IsAuthorized(user, plan, option string) (bool, error) {
	userRoles := make(map[string]bool)

	allowedPlans := make(map[string]bool)
	for _, role := range root.Users[user].Roles {
		for _, p := range root.Roles[role].AllowedPlans {
			allowedPlans[p] = true
		}
	}

	for p := range allowedPlans {
		if p == "*" {
			return true, nil
		}

		// TODO: is regex better?
		if strings.Contains(p, ".*") {
			for opt := range root.Plans[plan].Opts {
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
