package authorization

type Module interface {
	IsAuthorized(user, plan, option string) (bool, error)
}
