package authentication

import (
	"errors"

	"github.com/agrim123/gatekeeper/internal/pkg/access"
)

func Authenticate(email string) (bool, error) {
	if _, ok := access.Mappings[email]; !ok {
		return false, errors.New("Invalid user")
	}

	return true, nil
}
