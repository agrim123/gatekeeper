package actions

import "github.com/agrim123/gatekeeper/internal/pkg/remote"

type SSH struct {
	User       string `json:"user"`
	IP         string `json:"ip"`
	Port       string `json:"port"`
	PrivateKey string `json:"private_key"`
}

func (s SSH) Run(cmd string) error {
	remoteConnection := remote.NewRemoteConnection(s.User, s.IP, s.Port, s.PrivateKey)
	remoteConnection.MakeNewConnection()
	remoteConnection.RunCommand(cmd)
	remoteConnection.Close()
	return nil
}
