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

	User root.AccessMapping
}

func NewGatekeeper() *GateKeeper {
	return &GateKeeper{
		AuthenticationModule: authentication.NewDefaultModule(),
	}
}

func main() {
	if len(os.Args) < 2 {
		panic("Invalid arg")
	}

	config.Init()

	setup.Start()

	gatekeeper := NewGatekeeper()
	if username, authenticated, err := gatekeeper.AuthenticationModule.IsAuthenticated(); !authenticated {
		panic(err)
	} else {
		gatekeeper.User = root.Users[username]
	}

	plan := os.Args[1]

	if _, ok := root.Plans[plan]; !ok {
		panic("Invalid plan")
	}

	if len(os.Args) < 3 {
		fmt.Println(root.Plans[plan].Opts)
		return
	}

	option := os.Args[2]

	if _, ok := root.Plans[plan].Opts[option]; !ok {
		panic("Invalid option")
	}

	if authorized, err := authorization.IsAuthorized(gatekeeper.User.User.Email, plan, option); !authorized {
		panic(err)
	}

	// fmt.Println(root.Plans[plan].Opts[option].Run())
}

func (g *GateKeeper) SetAuthenticationModule(authenticationModule authentication.Module) {
	g.AuthenticationModule = authenticationModule
}
