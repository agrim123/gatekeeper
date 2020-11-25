package runtime

import (
	"context"
	"fmt"
	"log"

	"github.com/agrim123/gatekeeper/internal/constants"
	"github.com/agrim123/gatekeeper/internal/pkg/authorization"
	"github.com/agrim123/gatekeeper/internal/pkg/store"
)

type Runtime struct {
	ctx    context.Context
	plan   string
	option string

	AuthorizationModule authorization.Module
}

func NewDefaultRuntime() *Runtime {
	return &Runtime{
		AuthorizationModule: authorization.NewDefaultModule(),
	}
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

func (r *Runtime) authorize() {
	if authorized, err := r.AuthorizationModule.IsAuthorized(r.ctx, r.plan, r.option); !authorized {
		panic(err)
	} else {
		fmt.Println(fmt.Sprintf("Authorized `%s` to perform `%s %s`", r.ctx.Value(constants.UserContextKey).(string), r.plan, r.option))
	}
}

func (r *Runtime) verify() {
	if _, ok := store.Plans[r.plan]; !ok {
		allowedPlans := make([]string, 0)
		for plan := range store.Plans {
			allowedPlans = append(allowedPlans, plan)
		}

		log.Fatalf("Invalid plan: `%s`. Allowed plans: %v", r.plan, allowedPlans)
	}

	if r.option == "" {
		fmt.Println(store.Plans[r.plan].AllowedOptions())
		return
	}

	if _, ok := store.Plans[r.plan].Opts[r.option]; !ok {
		panic(fmt.Sprintf("Invalid option: %s. Allowed options: %v", r.option, store.Plans[r.plan].AllowedOptions()))
	}
}

func (r *Runtime) Prepare(ctx context.Context, plan, option string) {
	r.ctx = ctx
	r.setPlan(plan)
	r.setOption(option)
}

func (r *Runtime) Execute() error {
	r.verify()

	r.authorize()

	fmt.Println(fmt.Sprintf("Executing plan: %s %s", r.plan, r.option))
	return store.Plans[r.plan].Opts[r.option].Run()
}

func (r *Runtime) WithAuthorizationModule(authorizationModule authorization.Module) *Runtime {
	r.AuthorizationModule = authorizationModule
	return r
}
