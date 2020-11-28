package store

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/agrim123/gatekeeper/internal/pkg/containers"
)

type Option interface {
	Run() error
}

// Shell provides a shell on remote server
type Shell struct {
	Name   string
	Server string
}

func (s Shell) Run() error {
	server := Servers[s.Server]
	fmt.Println(server)

	cmd := exec.Command("ssh", "-i", server.Instances[0].PrivateKey, server.Instances[0].User+"@"+server.Instances[0].IP)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	return nil
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

// Local runs commands on local system
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

// Container runs command on remote server
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

	container.AddPreStage(*containers.NewStage([]string{
		"chmod",
		"400",
		"/home/deploy/keys/*",
	}, false))

	container.Create()
	container.Start(context.Background())
	container.TailLogs()
	container.Cleanup()

	return nil
}
