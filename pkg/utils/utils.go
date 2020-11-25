package utils

import (
	"context"
	"os"
	"os/user"

	"github.com/agrim123/gatekeeper/internal/constants"
)

func AttachExecutingUserToCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, constants.UserContextKey, getExecutingUser())
}

func getExecutingUser() string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	// TODO: [Fix] Not an enforcable check
	if os.Getenv("SUDO_USER") != "" {
		panic("Please run as non-sudo. Current real user: " + os.Getenv("SUDO_USER"))
		// return "", false, errors.New("Please run as non-sudo. Current real user: " + os.Getenv("SUDO_USER"))
	}

	return user.Username
}
