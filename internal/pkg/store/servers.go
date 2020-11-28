package store

import (
	"github.com/agrim123/gatekeeper/internal/constants"
	"github.com/agrim123/gatekeeper/internal/pkg/archive"
	"github.com/agrim123/gatekeeper/internal/pkg/filesystem"
	"github.com/agrim123/gatekeeper/pkg/services/remote"
)

type Instance struct {
	User       string `json:"user"`
	IP         string `json:"ip"`
	Port       string `json:"port"`
	PrivateKey string `json:"private_key"`
}

type Server struct {
	Name      string     `json:"name"`
	Instances []Instance `json:"instances"`
}

func (s Server) GetPrivateKeysTar() string {
	privateKeys := make([]string, 0)

	for _, instance := range s.Instances {
		privateKeys = append(privateKeys, instance.PrivateKey)
	}

	filesystem.CopyFilesToDir(privateKeys, constants.PrivateKeysStagingPath)

	archive.Tar(constants.PrivateKeysStagingPath, constants.RootStagingPath)

	return constants.PrivateKeysStagingTarPath
}

func (s Instance) Run(cmd string) error {
	remoteConnection := remote.NewRemoteConnection(s.User, s.IP, s.Port, s.PrivateKey)
	defer remoteConnection.Close()
	remoteConnection.MakeNewConnection()
	remoteConnection.RunCommand(cmd)

	return nil
}
