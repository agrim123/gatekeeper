package access

import (
	"encoding/json"

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

type Action struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Type        string                 `json:"type"`
	Command     string                 `json:"command"`
	Attributes  map[string]interface{} `json:"attributes"`
}

var Mappings map[string]AccessMapping
var Roles map[string]Role
var Actions map[string]Action

func Init() {
	var actions []Action
	actionsByteData, _ := json.Marshal(viper.Get("actions"))
	json.Unmarshal(actionsByteData, &actions)

	Actions = make(map[string]Action)
	for _, action := range actions {
		Actions[action.Name] = action
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
