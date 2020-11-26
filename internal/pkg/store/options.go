package store

import (
	"context"
	"fmt"
	"log"
	"os/exec"

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
	for _, stage := range l.Stages {
		fmt.Println("Running command: " + stage)
		out, err := exec.Command(stage).Output()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(out))
	}
	return nil
}

type Container struct {
	Name    string
	Server  string            `json:"server"`
	Stages  [][]string        `json:"stages"`
	Volumes map[string]string `json:"volumes"`
}

func (c Container) Run() error {
	container := containers.Container{
		Image:       "gatekeeper",
		Name:        "gatekeeper",
		Cmds:        c.Stages,
		Mounts:      c.Volumes,
		FilesToCopy: []string{Servers[c.Server].GetPrivateKeysTar()},
	}

	container.Create()
	container.Start(context.Background())
	container.TailLogs()
	container.Cleanup()

	return nil
}
