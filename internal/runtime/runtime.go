package runtime

import (
	"context"

	"github.com/agrim123/gatekeeper/pkg/logger"
	"github.com/agrim123/gatekeeper/pkg/store"
)

type Runtime struct {
	ctx context.Context
}

func NewDefaultRuntime() *Runtime {
	return &Runtime{}
}

func NewRuntime(ctx context.Context) *Runtime {
	r := NewDefaultRuntime()
	return r
}

func (r *Runtime) Execute(ctx context.Context, plan, option string) error {
	logger.Info("Executing plan: %s %s", plan, option)
	return store.Store.Plans[plan].Opts[option].Run()
}
