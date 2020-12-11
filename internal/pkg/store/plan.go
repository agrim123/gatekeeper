package store

import (
	"strings"
)

type Plan struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Options     map[string]interface{} `json:"options"`

	Opts map[string]Option `json:"-"`
}

func (p Plan) GetOptions() []string {
	options := make([]string, len(p.Opts))

	i := 0
	for opt := range p.Opts {
		options[i] = opt
		i++
	}

	return options
}

// GetAllowedCommands returns the allowed options for plans allowed for user
// returned format map[plan_name][array of options]
func GetAllowedCommands(user string) map[string][]string {
	cmds := make(map[string][]string)

	for _, role := range Users[user].Roles {
		for _, p := range Roles[role].AllowedPlans {
			options := []string{}
			if p == "*" {
				for pp := range Plans {
					for po := range Plans[pp].Opts {
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
				for po := range Plans[plan].Opts {
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

func GetAvailablePlans() []string {
	plans := make([]string, len(Plans))

	i := 0
	for plan := range Plans {
		plans[i] = plan
		i++
	}
	return plans
}
