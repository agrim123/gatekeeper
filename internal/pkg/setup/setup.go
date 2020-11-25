package setup

import (
	"encoding/json"

	"github.com/agrim123/gatekeeper/internal/pkg/store"
	"github.com/spf13/viper"
)

func Init() {
	var servers []store.Server
	serversByteData, _ := json.Marshal(viper.Get("servers"))
	json.Unmarshal(serversByteData, &servers)

	store.Servers = make(map[string]store.Server)
	for _, server := range servers {
		store.Servers[server.Name] = server
	}

	var plans []store.Plan
	plansByteData, _ := json.Marshal(viper.Get("plan"))
	json.Unmarshal(plansByteData, &plans)

	store.Plans = make(map[string]store.Plan)
	for _, plan := range plans {
		finalOptions := make(map[string]store.Option)
		for name, option := range plan.Options {
			if name == "deploy" {
				var deploy store.Remote
				deployBytesdata, _ := json.Marshal(option)
				json.Unmarshal(deployBytesdata, &deploy)
				finalOptions[name] = deploy
			} else if name == "status" {
				var status store.Remote
				statusBytesdata, _ := json.Marshal(option)
				json.Unmarshal(statusBytesdata, &status)
				finalOptions[name] = status
			}
		}

		plan.Opts = finalOptions
		plan.Options = nil

		store.Plans[plan.Name] = plan
	}

	var roles []store.Role
	rolesByteData, _ := json.Marshal(viper.Get("roles"))
	json.Unmarshal(rolesByteData, &roles)

	store.Roles = make(map[string]store.Role)
	for _, role := range roles {
		store.Roles[role.Name] = role
	}

	var users []store.AccessMapping
	usersByteData, _ := json.Marshal(viper.Get("users"))
	json.Unmarshal(usersByteData, &users)

	store.Users = make(map[string]store.AccessMapping)
	for _, user := range users {
		store.Users[user.User.Email] = user
	}
}
