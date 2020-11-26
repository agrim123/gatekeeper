package store

import (
	"context"
	"fmt"

	"github.com/agrim123/gatekeeper/internal/pkg/containers"
)

type Option interface {
	Run() error
}

type Remote struct {
	Name   string
	Server string   `json:"server"`
	Stages []string `json:"stages"`
}

func (r Remote) Run() error {
	server := Servers[r.Server]

	for _, instance := range server.Instances {
		// TODO: run mutiple command on same connection
		for _, command := range r.Stages {
			fmt.Println(command)
			instance.Run(command)
		}
	}

	return nil
}

type Local struct {
	Name   string
	Stages []string `json:"stages"`
}

func (l Local) Run() error {
	// for _, stage := range l.Stages {
	// 	fmt.Println("Running command: " + stage)
	// 	out, err := exec.Command(stage).Output()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	fmt.Println(string(out))
	// }

	c := containers.Container{
		Image: "gatekeeper",
		Name:  "gatekeeper-jail",
		Cmds:  []string{"/bin/ls", "-lh", "/keys"},
		Mounts: map[string]string{
			"<path>": "/keys",
		},
	}

	fmt.Println(c.Create())
	fmt.Println(c.Start(context.Background()))
	c.TailLogs()
	c.Cleanup()

	return nil
}
