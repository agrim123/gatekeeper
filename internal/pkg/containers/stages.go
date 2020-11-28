package containers

type Stage struct {
	user string `json:"-"`

	Privileged bool     `json:"privileged"`
	Command    []string `json:"command"`
}

func NewStage(command []string, privileged bool) *Stage {
	stage := &Stage{
		Command:    command,
		Privileged: privileged,
	}

	if privileged {
		stage.user = "root"
	} else {
		stage.user = "deploy"
	}

	return stage
}
