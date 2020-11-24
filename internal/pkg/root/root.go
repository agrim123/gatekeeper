package root

import "github.com/agrim123/gatekeeper/internal/pkg/remote"

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

type ServerOption struct {
	Server   string   `json:"server"`
	Commands []string `json:"commands"`
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
