package authorizaion

import "github.com/agrim123/gatekeeper/internal/pkg/access"

func IsAuthorizedToPerformAction(user *access.User, action string) bool {
	allowedActions := make(map[string]bool)
	for _, role := range access.Mappings[user.Email].Roles {
		roleActions := access.Roles[role]
		for _, action := range roleActions.Actions {
			allowedActions[action] = true
		}
	}

	_, ok := allowedActions[action]

	return ok
}
