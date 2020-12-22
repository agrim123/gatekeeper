package authentication

import (
	"context"
	"errors"

	"github.com/agrim123/gatekeeper/internal/pkg/utils"
	"github.com/agrim123/gatekeeper/internal/store"
	"github.com/agrim123/gatekeeper/pkg/logger"
)

type DefaultModule struct {
	ctx context.Context
}

func NewDefaultModule(ctx context.Context) *DefaultModule {
	return &DefaultModule{
		ctx: ctx,
	}
}

func (dm DefaultModule) IsAuthenticated() (bool, error) {
	username := utils.GetUsernameFromCtx(dm.ctx)

	if _, ok := store.Store.Users[username]; !ok {
		return false, errors.New("Invalid user: " + username)
	}

	logger.Success("Authenticated as " + logger.Underline(username))
	return true, nil
}
