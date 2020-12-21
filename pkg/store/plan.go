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
func (s *StoreStruct) GetAllowedCommands(user string) map[string][]string {
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
