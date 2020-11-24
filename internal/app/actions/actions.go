package actions

type Action interface {
	Run(cmd string) error
}

type ActionX struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Type        string                 `json:"type"`
	Command     string                 `json:"command"`
	Attributes  map[string]interface{} `json:"attributes"`

	Action Action `json:"-"`
}

func (a *ActionX) Run() {
	a.Action.Run(a.Command)
}
