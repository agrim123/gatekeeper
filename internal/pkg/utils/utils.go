package utils

import (
	"context"
	"os"
	"os/user"

	"github.com/agrim123/gatekeeper/internal/constants"
	"github.com/agrim123/gatekeeper/pkg/logger"
)

func getFromContext(ctx context.Context, key constants.ContextKeyType) string {
	keyInterface := ctx.Value(key)
	if keyInterface == nil {
		logger.Fatal("Unable to extract %v from context. Aborting.", key)
	}

	return keyInterface.(string)
}

func AttachExecutingUserToCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, constants.UserContextKey, GetExecutingUser())
}

func GetUsernameFromCtx(ctx context.Context) string {
	return getFromContext(ctx, constants.UserContextKey)
}

func GetExecutingUser() string {
	user, err := user.Current()
	if err != nil || user.Username == "" {
		panic(err)
	}

	// TODO: [Fix] Not an enforcable check
	if os.Getenv("SUDO_USER") != "" {
		logger.Fatal("Please run as non-sudo. Current real user: %s", logger.Bold(os.Getenv("SUDO_USER")))
	}

	return user.Username
}
