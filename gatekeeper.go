package gatekeeper

import (
	"context"
	"fmt"

	"github.com/agrim123/gatekeeper/internal/constants"
	"github.com/agrim123/gatekeeper/internal/guard"
	"github.com/agrim123/gatekeeper/internal/runtime"
	"github.com/agrim123/gatekeeper/pkg/authentication"
	"github.com/agrim123/gatekeeper/pkg/authorization"
	"github.com/agrim123/gatekeeper/pkg/filesystem"
	"github.com/agrim123/gatekeeper/pkg/notifier"
	"github.com/agrim123/gatekeeper/pkg/store"
)

type GateKeeper struct {
	ctx context.Context

	store *store.StoreStruct

	runtime *runtime.Runtime
	guard   *guard.Guard

	notifier notifier.Notifier
}

// NewGatekeeper returns new instance of gatekeeper with default modules
func NewGatekeeper(ctx context.Context, initStore *store.StoreStruct) *GateKeeper {
	// Initializes the staging path for containers
	filesystem.CreateDir(constants.RootStagingPath)

	store.Init(initStore)

	g := &GateKeeper{
		ctx:      ctx,
		runtime:  runtime.NewDefaultRuntime(),
		notifier: notifier.GetNotifier(),
		guard:    guard.NewGuard(),
		store:    store.Store,
	}

	return g
}

// WithNotifier updates the notifier module
func (g *GateKeeper) WithNotifier(notifier notifier.Notifier) *GateKeeper {
	g.notifier = notifier
	return g
}

// WithAuthorizationModule updates the guard's authorization module
func (g *GateKeeper) WithAuthorizationModule(authorizationModule authorization.Module) *GateKeeper {
	g.guard = g.guard.WithAuthorizationModule(authorizationModule)
	return g
}

// WithAuthenticationModule updates the guard's authentication module
func (g *GateKeeper) WithAuthenticationModule(authenticationModule authentication.Module) *GateKeeper {
	g.guard = g.guard.WithAuthenticationModule(authenticationModule)
	return g
}

// Run runs the command given to the gatekeeper
// It then delegates different tasks
func (g *GateKeeper) Run(plan, option string) {
	g.guard.Verify(g.ctx, plan, option)

	err := g.runtime.Execute(g.ctx, plan, option)
	if err != nil {
		g.notifier.Notify(
			fmt.Sprintf(
				"Plan `%s %s` executed by `%s` failed. Error: %s",
				plan,
				option,
				g.ctx.Value(constants.UserContextKey),
				err.Error(),
			),
		)
	} else {
		g.notifier.Notify(
			fmt.Sprintf(
				"Plan `%s %s` executed by `%s` successfully!",
				plan,
				option,
				g.ctx.Value(constants.UserContextKey),
			),
		)
	}
}
