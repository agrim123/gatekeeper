package gatekeeper

import (
	"context"

	"github.com/agrim123/gatekeeper/internal/constants"
	"github.com/agrim123/gatekeeper/internal/pkg/authentication"
	"github.com/agrim123/gatekeeper/internal/pkg/runtime"
	"github.com/agrim123/gatekeeper/internal/pkg/setup"
	"github.com/agrim123/gatekeeper/pkg/config"
)

type GateKeeper struct {
	ctx context.Context

	AuthenticationModule authentication.Module
	runtime              *runtime.Runtime
}

func NewGatekeeper(ctx context.Context) *GateKeeper {
	config.Init()

	setup.Init()

	return &GateKeeper{
		ctx:                  ctx,
		AuthenticationModule: authentication.NewDefaultModule(),
		runtime:              runtime.NewDefaultRuntime(),
	}
}

func (g *GateKeeper) WithAuthenticationModule(authenticationModule authentication.Module) *GateKeeper {
	g.AuthenticationModule = authenticationModule
	return g
}

func (g *GateKeeper) authenticate() {
	username, authenticated, err := g.AuthenticationModule.IsAuthenticated()
	if !authenticated {
		panic(err)
	} else {
		g.ctx = context.WithValue(g.ctx, constants.UserContextKey, username)
	}
}

func (g *GateKeeper) Run(plan, option string) {
	g.authenticate()

	g.runtime.Prepare(g.ctx, plan, option)

	g.runtime.Execute()
}
