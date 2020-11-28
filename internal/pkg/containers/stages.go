package containers

import "strings"

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
		stage.user = RootUser
	} else {
		stage.user = NonRootUser
	}

	return stage
}

func (s Stage) String() string {
	return strings.Join(s.Command, " ")
}
