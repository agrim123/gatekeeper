package store

import "fmt"

type Option interface {
	Run() error
}

type Remote struct {
	Server string   `json:"server"`
	Stages []string `json:"stages"`
}

func (r Remote) Run() error {
	server := Servers[r.Server]

	for _, instance := range server.Instances {
		for _, command := range r.Stages {
			fmt.Println(command)
			instance.Run(command)
		}
	}

	return nil
}

type Local struct {
	Stages []string `json:"stages"`
}

func (l Local) Run() error {
	return nil
}
