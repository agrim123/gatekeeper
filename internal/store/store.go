package store

import (
	"encoding/json"
	"strings"
)

var Store *StoreStruct

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

func (s *StoreStruct) WithServers(servers []Server) *StoreStruct {
	for _, server := range servers {
		server.NormalizeInstancesPrivateKeys()
		s.Servers[server.Name] = &server
	}

	return s
}

func (s *StoreStruct) WithPlans(plans []Plan) *StoreStruct {
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

	return s
}

func (s *StoreStruct) WithUsers(users []User) *StoreStruct {
	for _, user := range users {
		s.Users[user.User.Username] = user
	}

	return s
}

func (s *StoreStruct) WithGroups(groups []Group) *StoreStruct {
	for _, group := range groups {
		s.Groups[group.Name] = group
	}

	return s
}

// GetAllowedCommandsForUser returns the allowed options for plans allowed for user
// returned format map[plan_name][array of options]
func (s *StoreStruct) GetAllowedCommandsForUser(user string) map[string][]string {
	cmds := make(map[string][]string)

	for _, group := range s.Users[user].Groups {
		for _, p := range s.Groups[group].AllowedPlans {
			options := []string{}
			if p == "*" {
				for pp := range s.Plans {
					for po := range s.Plans[pp].Opts {
						options = append(options, po)
					}

					if allowedOptions, ok := cmds[pp]; ok {
						cmds[pp] = append(allowedOptions, options...)
					} else {
						cmds[pp] = options
					}
				}

				return cmds
			}

			allowedPlan := strings.Split(p, ".")
			if len(allowedPlan) != 2 {
				continue
			}

			plan := allowedPlan[0]
			option := allowedPlan[1]

			if option == "*" {
				for po := range s.Plans[plan].Opts {
					options = append(options, po)
				}
			} else {
				options = append(options, option)
			}

			if allowedOptions, ok := cmds[plan]; ok {
				cmds[plan] = append(allowedOptions, options...)
			} else {
				cmds[plan] = options
			}
		}
	}

	return cmds
}

func (s *StoreStruct) GetAvailablePlans() []string {
	plans := make([]string, len(s.Plans))

	i := 0
	for plan := range s.Plans {
		plans[i] = plan
		i++
	}

	return plans
}
