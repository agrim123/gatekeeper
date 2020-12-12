package store

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/agrim123/gatekeeper/internal/constants"
	"github.com/agrim123/gatekeeper/internal/pkg/containers"
	"github.com/agrim123/gatekeeper/pkg/logger"
	"github.com/agrim123/gatekeeper/pkg/services/remote"
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

	instance := server.Instances[0]
	if len(server.Instances) > 1 {
		choice := 0
		for {
			logger.Info("Multiple instances present in the cluster.")

			for index, ins := range server.Instances {
				logger.Info("[ %d ] %s", index, ins.String())
			}
			choice, _ = strconv.Atoi(logger.InfoScan("Choose which one to use: "))
			if choice < len(server.Instances) {
				instance = server.Instances[choice]
				break
			}

			logger.Error("Invalid choice %d", choice)
		}
	}

	logger.Info("Spawning shell for %s", logger.Bold(instance.String()))

	r := remote.NewRemoteConnection(instance.User, instance.IP, instance.Port, instance.PrivateKey)
	r.MakeNewConnection()
	r.SpawnShell()

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
		logger.Info("Running stages on %s", logger.Bold(instance.String()))
		instance.Run(r.Stages)
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
		Image:       constants.BaseImageName,
		Name:        constants.BaseContainerName,
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

	for _, instance := range Servers[c.Server].Instances {
		container.AddPreStage(*containers.NewStage(
			[]string{
				"chmod",
				"400",
				"/home/deploy/keys/" + filepath.Base(instance.PrivateKey),
			},
			false,
		))
	}

	container.Create()
	err := container.Start(context.Background())
	if err != nil {
		return errors.New("Unable to complete plan")
	}

	container.TailLogs()
	container.Cleanup()

	return nil
}
