package authentication

import (
	"context"
	"errors"

	"github.com/agrim123/gatekeeper/internal/constants"
	"github.com/agrim123/gatekeeper/internal/pkg/store"
	"github.com/agrim123/gatekeeper/pkg/logger"
)

type DefaultModule struct{}

func NewDefaultModule() *DefaultModule {
	return &DefaultModule{}
}

func (dm DefaultModule) IsAuthenticated(ctx context.Context) (bool, error) {
	username := ctx.Value(constants.UserContextKey).(string)

	if _, ok := store.Store.Users[username]; !ok {
		return false, errors.New("Invalid user: " + username)
	}

	logger.Success("Authenticated as " + logger.Underline(username))
	return true, nil
}
