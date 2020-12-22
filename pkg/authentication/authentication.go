package authentication

type Module interface {
	IsAuthenticated() (bool, error)
}
