package access

import (
	"encoding/json"

	"github.com/agrim123/gatekeeper/internal/app/actions"
	"github.com/spf13/viper"
)

type User struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type Role struct {
	Name    string   `json:"name"`
	Actions []string `json:"actions"`
}

type AccessMapping struct {
	User  User
	Roles []string
}

var Mappings map[string]AccessMapping
var Roles map[string]Role
var Actions map[string]actions.ActionX

func Init() {
	type actionTemp struct {
		actions.BaseAction
		Attributes map[string]interface{} `json:"attributes"`
	}

	var actionsInterface []actionTemp
	actionsByteData, _ := json.Marshal(viper.Get("actions"))
	json.Unmarshal(actionsByteData, &actionsInterface)

	Actions = make(map[string]actions.ActionX)
	for _, action := range actionsInterface {
		if action.Type == "ssh" {
			Actions[action.Name] = actions.ActionX{
				BaseAction: action.BaseAction,
				Attributes: actions.SSH{
					User:       action.Attributes["user"].(string),
					IP:         action.Attributes["ip"].(string),
					Port:       action.Attributes["port"].(string),
					PrivateKey: action.Attributes["private_key"].(string),
				},
			}
		}
	}

	var roles []Role
	rolesByteData, _ := json.Marshal(viper.Get("roles"))
	json.Unmarshal(rolesByteData, &roles)

	Roles = make(map[string]Role)
	for _, role := range roles {
		a := make([]string, 0)
		for _, action := range role.Actions {
			if _, ok := Actions[action]; ok {
				a = append(a, action)
			}
		}

		role.Actions = a

		Roles[role.Name] = role
	}

	var accessMappings []AccessMapping
	accessMappingsByteData, _ := json.Marshal(viper.Get("access_mappings"))
	json.Unmarshal(accessMappingsByteData, &accessMappings)

	Mappings = make(map[string]AccessMapping)
	for _, mapping := range accessMappings {
		Mappings[mapping.User.Email] = mapping
	}
}
