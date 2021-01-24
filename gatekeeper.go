package gatekeeper

import (
	"context"
	"fmt"

	"github.com/agrim123/gatekeeper/internal/constants"
	"github.com/agrim123/gatekeeper/internal/guard"
	"github.com/agrim123/gatekeeper/internal/runtime"
	"github.com/agrim123/gatekeeper/internal/store"
	"github.com/agrim123/gatekeeper/internal/utils"
	"github.com/agrim123/gatekeeper/pkg/authentication"
	"github.com/agrim123/gatekeeper/pkg/authorization"
	"github.com/agrim123/gatekeeper/pkg/logger"
	"github.com/agrim123/gatekeeper/pkg/notifier"
	storePkg "github.com/agrim123/gatekeeper/pkg/store"
)

type GateKeeper struct {
	ctx context.Context

	store *store.StoreStruct

	runtime *runtime.Runtime
	guard   *guard.Guard

	notifyRequester func(string)
}

// NewGatekeeper returns new instance of gatekeeper with default modules
func NewGatekeeper(ctx context.Context) *GateKeeper {
	// Initializes the staging path for containers
	// filesystem.CreateDir(constants.RootStagingPath)

	ctx = utils.AttachExecutingUserToCtx(ctx)

	g := &GateKeeper{
		ctx:             ctx,
		runtime:         runtime.NewRuntime(ctx),
		notifyRequester: notifier.AttachFallbackNotifier(notifier.NewDefaultNotifier()),
		guard:           guard.NewGuard(ctx),
		store:           store.Store,
	}

	return g
}

// WithNotifier updates the notifier module
func (g *GateKeeper) WithNotifier(customNotifier notifier.Notifier) *GateKeeper {
	g.notifyRequester = notifier.AttachFallbackNotifier(customNotifier)
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

func (g *GateKeeper) AllowedCommands() map[string][]string {
	return g.store.GetAllowedCommandsForUser(utils.GetExecutingUser())
}

func (g *GateKeeper) Whoami() *storePkg.Whoami {
	return &storePkg.Whoami{
		Username:        utils.GetExecutingUser(),
		Groups:          store.Store.Users[utils.GetExecutingUser()].Groups,
		AllowedCommands: g.AllowedCommands(),
	}
}

// Run runs the command given to the gatekeeper
// It then delegates different tasks
func (g *GateKeeper) Run(plan, option string) {
	g.guard.Verify(plan, option)

	err := g.runtime.Execute(plan, option)
	if err != nil {
		g.notifyRequester(
			fmt.Sprintf(
				"Plan `%s %s` executed by `%s` failed. Error: %s",
				logger.Underline(plan),
				logger.Underline(option),
				logger.Bold(g.ctx.Value(constants.UserContextKey).(string)),
				err.Error(),
			),
		)
	} else {
		g.notifyRequester(
			fmt.Sprintf(
				"Plan `%s %s` executed by `%s` successfully!",
				logger.Underline(plan),
				logger.Underline(option),
				logger.Bold(g.ctx.Value(constants.UserContextKey).(string)),
			),
		)
	}
}
