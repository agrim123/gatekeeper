package gatekeeper

import (
	"context"
	"fmt"

	"github.com/agrim123/gatekeeper/internal/constants"
	"github.com/agrim123/gatekeeper/internal/gatekeeper/runtime"
	"github.com/agrim123/gatekeeper/internal/pkg/authentication"
	"github.com/agrim123/gatekeeper/internal/pkg/notifier"
	"github.com/agrim123/gatekeeper/internal/pkg/setup"
	"github.com/agrim123/gatekeeper/pkg/config"
)

type GateKeeper struct {
	ctx context.Context

	AuthenticationModule authentication.Module
	runtime              *runtime.Runtime

	Notifier notifier.Notifier
}

func NewGatekeeper(ctx context.Context) *GateKeeper {
	config.Init()

	setup.Init()

	return &GateKeeper{
		ctx:                  ctx,
		AuthenticationModule: authentication.NewDefaultModule(),
		runtime:              runtime.NewDefaultRuntime(),
		Notifier:             notifier.GetNotifier(),
	}
}

func (g *GateKeeper) WithAuthenticationModule(authenticationModule authentication.Module) *GateKeeper {
	g.AuthenticationModule = authenticationModule
	return g
}

func (g *GateKeeper) WithNotifier(notifier notifier.Notifier) *GateKeeper {
	g.Notifier = notifier
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

	g.Notifier.Notify(fmt.Sprintf("Plan `%s %s` executed by `%s` successfully!", plan, option, g.ctx.Value(constants.UserContextKey)))
}
