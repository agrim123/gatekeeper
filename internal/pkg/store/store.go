package store

import (
	"github.com/agrim123/gatekeeper/internal/constants"
	"github.com/agrim123/gatekeeper/internal/pkg/archive"
	"github.com/agrim123/gatekeeper/internal/pkg/filesystem"
	"github.com/agrim123/gatekeeper/pkg/services/remote"
)

type User struct {
	Email string `json:"email"`
}

type AccessMapping struct {
	User  User
	Roles []string
}

type Role struct {
	Name         string   `json:"name"`
	AllowedPlans []string `json:"allowed_plans"`
}

type Plan struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Options     map[string]interface{} `json:"options"`

	Opts map[string]Option `json:"-"`
}

type Instance struct {
	User       string `json:"user"`
	IP         string `json:"ip"`
	Port       string `json:"port"`
	PrivateKey string `json:"private_key"`
}

func (s Instance) Run(cmd string) error {
	remoteConnection := remote.NewRemoteConnection(s.User, s.IP, s.Port, s.PrivateKey)
	remoteConnection.MakeNewConnection()
	remoteConnection.RunCommand(cmd)
	remoteConnection.Close()
	return nil
}

type Server struct {
	Name      string     `json:"name"`
	Instances []Instance `json:"instances"`
}

var Users map[string]AccessMapping
var Roles map[string]Role
var Servers map[string]Server
var Plans map[string]Plan

func (p Plan) AllowedOptions() []string {
	options := make([]string, len(p.Opts))

	i := 0
	for opt := range p.Opts {
		options[i] = opt
		i++
	}

	return options
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
