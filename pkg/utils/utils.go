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
	return context.WithValue(ctx, constants.UserContextKey, getExecutingUser())
}

func AttachPlanToCtx(ctx context.Context, plan string) context.Context {
	return context.WithValue(ctx, constants.PlanContextKey, plan)
}

func AttachOptionToCtx(ctx context.Context, option string) context.Context {
	return context.WithValue(ctx, constants.OptionContextKey, option)
}

func GetUsernameFromCtx(ctx context.Context) string {
	return getFromContext(ctx, constants.UserContextKey)
}

func GetPlanFromCtx(ctx context.Context) string {
	return getFromContext(ctx, constants.PlanContextKey)
}

func GetOptionFromCtx(ctx context.Context) string {
	return getFromContext(ctx, constants.OptionContextKey)
}

func getExecutingUser() string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	// TODO: [Fix] Not an enforcable check
	if os.Getenv("SUDO_USER") != "" {
		logger.Fatal("Please run as non-sudo. Current real user: %s", logger.Bold(os.Getenv("SUDO_USER")))
	}

	return user.Username
}
