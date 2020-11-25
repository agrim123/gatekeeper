package authentication

import (
	"errors"
	"fmt"
	"os"
	"os/user"

	"github.com/agrim123/gatekeeper/internal/pkg/store"
)

type DefaultModule struct{}

func NewDefaultModule() *DefaultModule {
	return &DefaultModule{}
}

func (dm DefaultModule) IsAuthenticated() (string, bool, error) {
	user, err := user.Current()
	if err != nil {
		return "", false, err
	}

	if _, ok := store.Users[user.Username]; !ok {
		return "", false, errors.New("Invalid user: " + user.Username)
	}

	// TODO: [Fix] Not an enforcable check
	if os.Getenv("SUDO_USER") != "" {
		return "", false, errors.New("Please run as non-sudo. Current real user: " + os.Getenv("SUDO_USER"))
	}

	fmt.Println("Authenticated as " + user.Username)
	return user.Username, true, nil
}
