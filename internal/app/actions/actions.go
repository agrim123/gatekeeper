package actions

type Action interface {
	Run(cmd string) error
}

type BaseAction struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Command     string `json:"command"`
}

type ActionX struct {
	BaseAction

	Action Action
}

func (a *ActionX) Run() {
	a.Action.Run(a.Command)
}
