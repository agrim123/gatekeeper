package authorization

type Module interface {
	IsAuthorized(plan, option string) (bool, error)
}
