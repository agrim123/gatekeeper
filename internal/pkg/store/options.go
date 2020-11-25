package store

import (
	"fmt"
	"log"
	"os/exec"
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
