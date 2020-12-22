package runtime

import (
	"context"

	"github.com/agrim123/gatekeeper/internal/store"
	"github.com/agrim123/gatekeeper/pkg/logger"
)

type Runtime struct {
	ctx context.Context
}

func NewRuntime(ctx context.Context) *Runtime {
	return &Runtime{
		ctx: ctx,
	}
}

func (r *Runtime) Execute(plan, option string) error {
	logger.Info("Executing plan: %s %s", plan, option)
	return store.Store.Plans[plan].Opts[option].Run()
}
