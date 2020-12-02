package runtime

import (
	"context"
	"fmt"

	"github.com/agrim123/gatekeeper/internal/pkg/store"
	"github.com/agrim123/gatekeeper/pkg/logger"
)

type Runtime struct {
	ctx    context.Context
	plan   string
	option string
}

func NewDefaultRuntime() *Runtime {
	return &Runtime{}
}

func NewRuntime(ctx context.Context, plan, option string) *Runtime {
	r := NewDefaultRuntime()
	r.plan = plan
	r.option = option
	return r
}

func (r *Runtime) setPlan(plan string) {
	r.plan = plan
}

func (r *Runtime) setOption(option string) {
	r.option = option
}

func (r *Runtime) verify() {
	if _, ok := store.Plans[r.plan]; !ok {
		allowedPlans := make([]string, 0)
		for plan := range store.Plans {
			allowedPlans = append(allowedPlans, plan)
		}

		logger.Fatalf("Invalid plan: `%s`. Allowed plans: %v", r.plan, allowedPlans)
	}

	if r.option == "" {
		fmt.Println(store.Plans[r.plan].AllowedOptions())
		return
	}

	if _, ok := store.Plans[r.plan].Opts[r.option]; !ok {
		logger.Fatalf("Invalid option: %s. Allowed options: %v", r.option, store.Plans[r.plan].AllowedOptions())
	}
}

func (r *Runtime) Prepare(ctx context.Context, plan, option string) {
	r.ctx = ctx
	r.setPlan(plan)
	r.setOption(option)
}

func (r *Runtime) Execute() error {
	r.verify()

	logger.Infof("Executing plan: %s %s", r.plan, r.option)
	return store.Plans[r.plan].Opts[r.option].Run()
}
