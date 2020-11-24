package main

import (
	"github.com/agrim123/gatekeeper/internal/pkg/access"
	"github.com/agrim123/gatekeeper/pkg/config"
)

func main() {
	config.Init()
	access.Init()
	// fmt.Println(authorizaion.IsAuthorizedToPerformAction(&access.User{Email: "agrim@xyz.com"}, "action1"))
	ssh := access.Actions["action1"]

	ssh.Run()
}
