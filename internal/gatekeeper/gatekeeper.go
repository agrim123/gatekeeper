package gatekeeper

import (
	"context"
	"fmt"

	"github.com/agrim123/gatekeeper/internal/constants"
	"github.com/agrim123/gatekeeper/internal/gatekeeper/guard"
	"github.com/agrim123/gatekeeper/internal/gatekeeper/runtime"
	"github.com/agrim123/gatekeeper/internal/pkg/filesystem"
	"github.com/agrim123/gatekeeper/internal/pkg/notifier"
)

type GateKeeper struct {
	ctx context.Context

	runtime *runtime.Runtime
	guard   *guard.Guard

	notifier notifier.Notifier
}

func NewGatekeeper(ctx context.Context) *GateKeeper {
	filesystem.CreateDir(constants.RootStagingPath)

	return &GateKeeper{
		ctx:      ctx,
		runtime:  runtime.NewDefaultRuntime(),
		notifier: notifier.GetNotifier(),
		guard:    guard.NewGuard(),
	}
}

func (g *GateKeeper) WithNotifier(notifier notifier.Notifier) *GateKeeper {
	g.notifier = notifier
	return g
}

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
