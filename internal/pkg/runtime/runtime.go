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
	Username string

	Plan   string
	Option string

	AuthorizationModule authorization.Module
}

func NewDefaultRuntime() *Runtime {
	return &Runtime{
		AuthorizationModule: authorization.NewDefaultModule(),
	}
}

func NewRuntime(ctx context.Context, plan, option string) *Runtime {
	r := NewDefaultRuntime()
	r.Plan = plan
	r.Option = option
	return r
}

func (r *Runtime) SetPlan(plan string) {
	r.Plan = plan
}

func (r *Runtime) SetOption(option string) {
	r.Option = option
}

func (r *Runtime) authorize() {
	if authorized, err := r.AuthorizationModule.IsAuthorized(r.Username, r.Plan, r.Option); !authorized {
		panic(err)
	} else {
		fmt.Println(fmt.Sprintf("Authorized `%s` to perform `%s %s`", r.Username, r.Plan, r.Option))
	}
}

func (r *Runtime) verify() {
	if _, ok := store.Plans[r.Plan]; !ok {
		allowedPlans := make([]string, 0)
		for plan := range store.Plans {
			allowedPlans = append(allowedPlans, plan)
		}

		log.Fatalf("Invalid plan: `%s`. Allowed plans: %v", r.Plan, allowedPlans)
	}

	if r.Option == "" {
		fmt.Println(store.Plans[r.Plan].AllowedOptions())
		return
	}

	if _, ok := store.Plans[r.Plan].Opts[r.Option]; !ok {
		panic(fmt.Sprintf("Invalid option: %s. Allowed options: %v", r.Option, store.Plans[r.Plan].AllowedOptions()))
	}
}

func (r *Runtime) Prepare(ctx context.Context, plan, option string) {
	r.Username = ctx.Value(constants.UserContextKey).(string)
	r.Plan = plan
	r.Option = option
}

func (r *Runtime) Execute() {
	r.verify()

	r.authorize()

	fmt.Println(fmt.Sprintf("Executing plan: %s %s", r.Plan, r.Option))
	// fmt.Println(store.Plans[plan].Opts[option].Run())
}

func (r *Runtime) WithAuthorizationModule(authorizationModule authorization.Module) *Runtime {
	r.AuthorizationModule = authorizationModule
	return r
}
