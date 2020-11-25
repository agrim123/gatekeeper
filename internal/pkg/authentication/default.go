package authentication

import (
	"context"
	"errors"
	"fmt"

	"github.com/agrim123/gatekeeper/internal/constants"
	"github.com/agrim123/gatekeeper/internal/pkg/store"
)

type DefaultModule struct{}

func NewDefaultModule() *DefaultModule {
	return &DefaultModule{}
}

func (dm DefaultModule) IsAuthenticated(ctx context.Context) (bool, error) {
	username := ctx.Value(constants.UserContextKey).(string)

	if _, ok := store.Users[username]; !ok {
		return false, errors.New("Invalid user: " + username)
	}

	fmt.Println("Authenticated as " + username)
	return true, nil
}
