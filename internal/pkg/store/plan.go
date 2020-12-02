package store

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
func GetAvailablePlans() []string {
	plans := make([]string, len(Plans))

	i := 0
	for plan := range Plans {
		plans[i] = plan
		i++
	}
	return plans
}
