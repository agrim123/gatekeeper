package main

import (
	"fmt"
	"os"

	"github.com/agrim123/gatekeeper/internal/pkg/authentication"
	"github.com/agrim123/gatekeeper/internal/pkg/authorization"
	"github.com/agrim123/gatekeeper/internal/pkg/root"
	"github.com/agrim123/gatekeeper/internal/pkg/setup"
	"github.com/agrim123/gatekeeper/pkg/config"
)

type GateKeeper struct {
	AuthenticationModule authentication.Module
	AuthorizationModule  authorization.Module

	User root.AccessMapping
}

func NewGatekeeper() *GateKeeper {
	return &GateKeeper{
		AuthenticationModule: authentication.NewDefaultModule(),
		AuthorizationModule:  authorization.NewDefaultModule(),
	}
}

func main() {
	if len(os.Args) < 2 {
		panic("Invalid arg")
	}

	config.Init()

	setup.Start()

	gatekeeper := NewGatekeeper()
	username, authenticated, err := gatekeeper.AuthenticationModule.IsAuthenticated()
	if !authenticated {
		panic(err)
	} else {
		gatekeeper.User = root.Users[username]
	}

	plan := os.Args[1]

	if _, ok := root.Plans[plan]; !ok {
		panic("Invalid plan")
	}

	if len(os.Args) < 3 {
		fmt.Println(root.Plans[plan].AllowedOptions())
		return
	}

	option := os.Args[2]

	if _, ok := root.Plans[plan].Opts[option]; !ok {
		panic(fmt.Sprintf("Invalid option: %s. Allowed options: %v", option, root.Plans[plan].AllowedOptions()))
	}

	if authorized, err := gatekeeper.AuthorizationModule.IsAuthorized(gatekeeper.User.User.Email, plan, option); !authorized {
		panic(err)
	} else {
		fmt.Println(fmt.Sprintf("Authorized `%s` to perform `%s %s`", username, plan, option))
	}

	// fmt.Println(root.Plans[plan].Opts[option].Run())
}

func (g *GateKeeper) SetAuthenticationModule(authenticationModule authentication.Module) {
	g.AuthenticationModule = authenticationModule
}

func (g *GateKeeper) SetAuthorizationModule(authorizationModule authorization.Module) {
	g.AuthorizationModule = authorizationModule
}
