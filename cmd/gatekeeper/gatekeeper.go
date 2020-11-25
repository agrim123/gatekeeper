package gatekeeper

import (
	"fmt"

	"github.com/agrim123/gatekeeper/internal/pkg/authentication"
	"github.com/agrim123/gatekeeper/internal/pkg/authorization"
	"github.com/agrim123/gatekeeper/internal/pkg/runtime"
	"github.com/agrim123/gatekeeper/internal/pkg/setup"
	"github.com/agrim123/gatekeeper/internal/pkg/store"
	"github.com/agrim123/gatekeeper/pkg/config"
)

type GateKeeper struct {
	AuthenticationModule authentication.Module
	AuthorizationModule  authorization.Module

	User store.AccessMapping

	Runtime *runtime.Runtime
}

func NewGatekeeper() *GateKeeper {
	return &GateKeeper{
		AuthenticationModule: authentication.NewDefaultModule(),
		AuthorizationModule:  authorization.NewDefaultModule(),
		Runtime:              runtime.NewDefaultRuntime(),
	}
}

func (g *GateKeeper) WithRuntime(runtime *runtime.Runtime) *GateKeeper {
	g.Runtime = runtime
	return g
}

func (g *GateKeeper) WithAuthenticationModule(authenticationModule authentication.Module) *GateKeeper {
	g.AuthenticationModule = authenticationModule
	return g
}

func (g *GateKeeper) WithAuthorizationModule(authorizationModule authorization.Module) *GateKeeper {
	g.AuthorizationModule = authorizationModule
	return g
}

func (g *GateKeeper) init() {
	config.Init()

	setup.Init()
}

func (g *GateKeeper) runPrechecks() {
	g.init()

	username, authenticated, err := g.AuthenticationModule.IsAuthenticated()
	if !authenticated {
		panic(err)
	} else {
		g.User = store.Users[username]
	}
}

func (g *GateKeeper) Run() {
	g.runPrechecks()

	g.Runtime.Verify()

	if authorized, err := g.AuthorizationModule.IsAuthorized(g.User.User.Email, g.Runtime.Plan, g.Runtime.Option); !authorized {
		panic(err)
	} else {
		fmt.Println(fmt.Sprintf("Authorized `%s` to perform `%s %s`", g.User.User.Email, g.Runtime.Plan, g.Runtime.Option))
	}

	g.Runtime.Execute()
}
