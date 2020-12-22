package guard

import (
	"context"

	"github.com/agrim123/gatekeeper/internal/constants"
	"github.com/agrim123/gatekeeper/pkg/authentication"
	"github.com/agrim123/gatekeeper/pkg/authorization"
	"github.com/agrim123/gatekeeper/pkg/logger"
)

type Guard struct {
	ctx context.Context

	authenticationModule authentication.Module
	authorizationModule  authorization.Module
}

func NewGuard(ctx context.Context) *Guard {
	return &Guard{
		ctx:                  ctx,
		authenticationModule: authentication.NewDefaultModule(ctx),
		authorizationModule:  authorization.NewDefaultModule(ctx),
	}
}

func (g *Guard) WithAuthorizationModule(authorizationModule authorization.Module) *Guard {
	g.authorizationModule = authorizationModule
	return g
}

func (g *Guard) WithAuthenticationModule(authenticationModule authentication.Module) *Guard {
	g.authenticationModule = authenticationModule
	return g
}

func (g *Guard) Verify(plan, option string) {
	g.authenticate()
	g.authorize(plan, option)
}

func (g *Guard) authenticate() {
	if authenticated, err := g.authenticationModule.IsAuthenticated(); !authenticated {
		logger.Fatal(err.Error())
	}
}

func (g *Guard) authorize(plan, option string) {
	if authorized, err := g.authorizationModule.IsAuthorized(plan, option); !authorized {
		logger.Fatal(err.Error())
	}

	logger.Success("Authorized `%s` to perform `%s %s`", logger.Underline(g.ctx.Value(constants.UserContextKey).(string)), plan, option)
}
