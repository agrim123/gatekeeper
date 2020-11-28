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
	Name      string
	Server    string             `json:"server"`
	Stages    []containers.Stage `json:"stages"`
	Volumes   map[string]string  `json:"volumes"`
	Protected bool               `json:"protected"`
}

func (c Container) Run() error {
	stages := make([]containers.Stage, len(c.Stages))
	for i, stage := range c.Stages {
		stages[i] = *containers.NewStage(stage.Command, stage.Privileged)
	}

	container := containers.Container{
		Image:       "gatekeeper",
		Name:        "gatekeeper",
		Stages:      stages,
		Mounts:      c.Volumes,
		FilesToCopy: []string{Servers[c.Server].GetPrivateKeysTar()},
		Protected:   c.Protected,
	}

	container.AddPreStage(*containers.NewStage([]string{
		"chown",
		"-R",
		"deploy:deploy",
		"/home/deploy/keys",
	}, true))

	container.Create()
	container.Start(context.Background())
	container.TailLogs()
	container.Cleanup()

	return nil
}
