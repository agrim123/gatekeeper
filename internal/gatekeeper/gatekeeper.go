package gatekeeper

import (
	"context"
	"fmt"

	"github.com/agrim123/gatekeeper/internal/constants"
	"github.com/agrim123/gatekeeper/internal/gatekeeper/guard"
	"github.com/agrim123/gatekeeper/internal/gatekeeper/runtime"
	"github.com/agrim123/gatekeeper/internal/pkg/filesystem"
	"github.com/agrim123/gatekeeper/internal/pkg/notifier"
	"github.com/agrim123/gatekeeper/internal/pkg/setup"
	"github.com/agrim123/gatekeeper/pkg/config"
)

type GateKeeper struct {
	ctx context.Context

	runtime *runtime.Runtime
	guard   *guard.Guard

	Notifier notifier.Notifier
}

func NewGatekeeper(ctx context.Context) *GateKeeper {
	config.Init()

	setup.Init()

	filesystem.CreateDir(constants.RootStagingPath)

	return &GateKeeper{
		ctx:      ctx,
		runtime:  runtime.NewDefaultRuntime(),
		Notifier: notifier.GetNotifier(),
		guard:    guard.NewGuard(),
	}
}

func (g *GateKeeper) WithNotifier(notifier notifier.Notifier) *GateKeeper {
	g.Notifier = notifier
	return g
}

func (g *GateKeeper) Run(plan, option string) {
	g.guard.Verify(g.ctx, plan, option)

	g.runtime.Prepare(g.ctx, plan, option)

	err := g.runtime.Execute()
	if err != nil {
		g.Notifier.Notify(
			fmt.Sprintf(
				"Plan `%s %s` executed by `%s` failed. Error: %s",
				plan,
				option,
				g.ctx.Value(constants.UserContextKey),
				err.Error(),
			),
		)
	} else {
		g.Notifier.Notify(
			fmt.Sprintf(
				"Plan `%s %s` executed by `%s` successfully!",
				plan,
				option,
				g.ctx.Value(constants.UserContextKey),
			),
		)
	}
}
