package store

import (
	"strings"

	"github.com/agrim123/gatekeeper/internal/constants"
	"github.com/agrim123/gatekeeper/internal/pkg/filesystem"
	"github.com/agrim123/gatekeeper/internal/pkg/filesystem/archive"
	"github.com/agrim123/gatekeeper/internal/pkg/services/remote"
)

type Instance struct {
	User       string `json:"user"`
	IP         string `json:"ip"`
	Port       string `json:"port"`
	PrivateKey string `json:"private_key"`
}

func (i *Instance) String() string {
	return i.User + "@" + i.IP
}

type Server struct {
	Name      string      `json:"name"`
	Instances []*Instance `json:"instances"`
}

func (s *Server) GetPrivateKeysTar() string {
	privateKeys := make([]string, 0)

	for _, instance := range s.Instances {
		privateKeys = append(privateKeys, instance.PrivateKey)
	}

	filesystem.CopyFilesToDir(privateKeys, constants.PrivateKeysStagingPath)

	archive.Tar(constants.PrivateKeysStagingPath, constants.RootStagingPath)

	return constants.PrivateKeysStagingTarPath
}

func (s *Server) NormalizeInstancesPrivateKeys() {
	instances := make([]*Instance, len(s.Instances))

	for i, instance := range s.Instances {
		instances[i] = instance
	}

	s.Instances = instances
}

func (i *Instance) Run(cmds []string) error {
	remoteConnection, err := remote.NewRemoteConnection(i.User, i.IP, i.Port, i.PrivateKey)
	defer remoteConnection.Close()
	if err != nil {
		return err
	}

	err = remoteConnection.MakeNewConnection()
	if err != nil {
		return err
	}

	return remoteConnection.RunCommand(strings.Join(cmds, "; "))
}
