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

type Attributes Action

type ActionX struct {
	BaseAction

	Attributes
}

func (a *ActionX) Run() {
	a.Attributes.Run(a.Command)
}
