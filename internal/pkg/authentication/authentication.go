package authentication

import (
	"errors"
	"fmt"
	"os"
	"os/user"

	"github.com/agrim123/gatekeeper/internal/pkg/root"
)

type Module interface {
	IsAuthenticated() (string, bool, error)
}

type DefaultModule struct{}

func (dm DefaultModule) IsAuthenticated() (string, bool, error) {
	user, err := user.Current()
	if err != nil {
		return "", false, err
	}

	if _, ok := root.Users[user.Username]; !ok {
		return "", false, errors.New("Invalid user")
	}

	// TODO: [Fix] Not an enforcable check
	if os.Getenv("SUDO_USER") != "" {
		return "", false, errors.New("Please run as non-sudo. Real user: " + os.Getenv("SUDO_USER"))
	}

	fmt.Println("Authenticated as " + user.Username)
	return user.Username, true, nil
}

func NewDefaultModule() *DefaultModule {
	return &DefaultModule{}
}
