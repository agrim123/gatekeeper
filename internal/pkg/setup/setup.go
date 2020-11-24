package setup

import (
	"encoding/json"

	"github.com/agrim123/gatekeeper/internal/pkg/root"
	"github.com/spf13/viper"
)

func Start() {
	var servers []root.Server
	serversByteData, _ := json.Marshal(viper.Get("servers"))
	json.Unmarshal(serversByteData, &servers)

	root.Servers = make(map[string]root.Server)
	for _, server := range servers {
		root.Servers[server.Name] = server
	}

	var plans []root.Plan
	plansByteData, _ := json.Marshal(viper.Get("plan"))
	json.Unmarshal(plansByteData, &plans)

	root.Plans = make(map[string]root.Plan)
	for _, plan := range plans {
		finalOptions := make(map[string]root.Option)
		for name, option := range plan.Options {
			if name == "deploy" {
				var deploy root.Deploy
				deployBytesdata, _ := json.Marshal(option)
				json.Unmarshal(deployBytesdata, &deploy)
				finalOptions[name] = deploy
			}
		}

		plan.Opts = finalOptions
		plan.Options = nil

		root.Plans[plan.Name] = plan
	}

	var roles []root.Role
	rolesByteData, _ := json.Marshal(viper.Get("roles"))
	json.Unmarshal(rolesByteData, &roles)

	root.Roles = make(map[string]root.Role)
	for _, role := range roles {
		root.Roles[role.Name] = role
	}

	var users []root.AccessMapping
	usersByteData, _ := json.Marshal(viper.Get("users"))
	json.Unmarshal(usersByteData, &users)

	root.Users = make(map[string]root.AccessMapping)
	for _, user := range users {
		root.Users[user.User.Email] = user
	}
}
