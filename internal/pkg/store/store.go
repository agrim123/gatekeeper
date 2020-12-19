package store

import (
	"encoding/json"

	"github.com/spf13/viper"
)

var Store StoreStruct

type StoreStruct struct {
	Users   map[string]User
	Groups  map[string]Group
	Servers map[string]*Server
	Plans   map[string]Plan
}

func NewStore() *StoreStruct {
	return &StoreStruct{
		Users:   make(map[string]User),
		Groups:  make(map[string]Group),
		Servers: make(map[string]*Server),
		Plans:   make(map[string]Plan),
	}
}

func Init() {
	s := NewStore()
	s.Initialize()
	Store = *s
}

func (s *StoreStruct) Initialize() {
	var groups []Group
	groupsByteData, _ := json.Marshal(viper.Get("groups"))
	json.Unmarshal(groupsByteData, &groups)
	s.PopulateGroups(groups)

	var users []User
	usersByteData, _ := json.Marshal(viper.Get("users"))
	json.Unmarshal(usersByteData, &users)
	s.PopulateUsers(users)

	var servers []Server
	serversByteData, _ := json.Marshal(viper.Get("servers"))
	json.Unmarshal(serversByteData, &servers)
	s.PopulateServers(servers)

	var plans []Plan
	plansByteData, _ := json.Marshal(viper.Get("plan"))
	json.Unmarshal(plansByteData, &plans)
	s.PopulatePlans(plans)
}

func (s *StoreStruct) PopulateServers(servers []Server) {
	for _, server := range servers {
		server.NormalizeInstancesPrivateKeys()
		s.Servers[server.Name] = &server
	}
}
func (s *StoreStruct) PopulatePlans(plans []Plan) {
	for _, plan := range plans {
		finalOptions := make(map[string]Option)
		for name, optionInterface := range plan.Options {
			option := optionInterface.(map[string]interface{})

			switch option["type"].(string) {
			case "remote":
				var remote Remote
				remoteBytesdata, _ := json.Marshal(option)
				json.Unmarshal(remoteBytesdata, &remote)
				finalOptions[name] = remote
			case "local":
				var local Local
				localBytesdata, _ := json.Marshal(option)
				json.Unmarshal(localBytesdata, &local)
				finalOptions[name] = local
			case "container":
				continue
				var container Container
				containerBytesdata, _ := json.Marshal(option)
				json.Unmarshal(containerBytesdata, &container)
				finalOptions[name] = container
			case "shell":
				var shell Shell
				shellBytesdata, _ := json.Marshal(option)
				json.Unmarshal(shellBytesdata, &shell)
				finalOptions[name] = shell
			}
		}

		plan.Opts = finalOptions
		plan.Options = nil

		s.Plans[plan.Name] = plan
	}
}

func (s *StoreStruct) PopulateUsers(users []User) {
	for _, user := range users {
		s.Users[user.User.Username] = user
	}
}

func (s *StoreStruct) PopulateGroups(groups []Group) {
	for _, group := range groups {
		s.Groups[group.Name] = group
	}
}
