package runtime

import (
	"fmt"
	"log"

	"github.com/agrim123/gatekeeper/internal/pkg/store"
)

type Runtime struct {
	Plan   string
	Option string
}

func NewDefaultRuntime() *Runtime {
	return &Runtime{}
}

func NewRuntime(plan, option string) *Runtime {
	return &Runtime{
		Plan:   plan,
		Option: option,
	}
}

func (r *Runtime) SetPlan(plan string) {
	r.Plan = plan
}

func (r *Runtime) SetOption(option string) {
	r.Option = option
}

func (r *Runtime) Verify() {
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

func (r *Runtime) Execute() {
	fmt.Println(fmt.Sprintf("Executing plan: %s %s", r.Plan, r.Option))
	// fmt.Println(store.Plans[plan].Opts[option].Run())
}
