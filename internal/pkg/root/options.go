package root

import "fmt"

type Option interface {
	Run() error
}

type Deploy struct {
	Server   string   `json:"server"`
	Commands []string `json:"commands"`
}

func (d Deploy) Run() error {
	server := Servers[d.Server]

	for _, instance := range server.Instances {
		for _, command := range d.Commands {
			fmt.Println(command)
			instance.Run(command)
		}
	}

	return nil
}
