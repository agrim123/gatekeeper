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

	store.Servers = make(map[string]*store.Server)
	for _, server := range servers {
		server.NormalizeInstancesPrivateKeys()
		store.Servers[server.Name] = &server
	}

	var plans []store.Plan
	plansByteData, _ := json.Marshal(viper.Get("plan"))
	json.Unmarshal(plansByteData, &plans)

	store.Plans = make(map[string]store.Plan)
	for _, plan := range plans {
		finalOptions := make(map[string]store.Option)
		for name, optionInterface := range plan.Options {
			option := optionInterface.(map[string]interface{})

			switch option["type"].(string) {
			case "remote":
				var remote store.Remote
				remoteBytesdata, _ := json.Marshal(option)
				json.Unmarshal(remoteBytesdata, &remote)
				finalOptions[name] = remote
			case "local":
				var local store.Local
				localBytesdata, _ := json.Marshal(option)
				json.Unmarshal(localBytesdata, &local)
				finalOptions[name] = local
			case "container":
				continue
				var container store.Container
				containerBytesdata, _ := json.Marshal(option)
				json.Unmarshal(containerBytesdata, &container)
				finalOptions[name] = container
			case "shell":
				var shell store.Shell
				shellBytesdata, _ := json.Marshal(option)
				json.Unmarshal(shellBytesdata, &shell)
				finalOptions[name] = shell
			}
		}

		plan.Opts = finalOptions
		plan.Options = nil

		store.Plans[plan.Name] = plan
	}

	var groups []store.Group
	groupsByteData, _ := json.Marshal(viper.Get("groups"))
	json.Unmarshal(groupsByteData, &groups)

	store.Groups = make(map[string]store.Group)
	for _, group := range groups {
		store.Groups[group.Name] = group
	}

	var users []store.User
	usersByteData, _ := json.Marshal(viper.Get("users"))
	json.Unmarshal(usersByteData, &users)

	store.Users = make(map[string]store.User)
	for _, user := range users {
		store.Users[user.User.Username] = user
	}
}
