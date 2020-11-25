package authentication

type Module interface {
	IsAuthenticated() (string, bool, error)
}
