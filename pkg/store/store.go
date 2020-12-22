package store

import (
	"encoding/json"

	"github.com/agrim123/gatekeeper/internal/pkg/utils"
	"github.com/agrim123/gatekeeper/internal/store"
	"github.com/agrim123/gatekeeper/pkg/logger"
)

func InitStore(users, plan, servers, groups interface{}) {
	var groupsStruct []store.Group
	groupsByteData, _ := json.Marshal(groups)
	err := json.Unmarshal(groupsByteData, &groupsStruct)
	if err != nil {
		logger.Fatal("Unable to parse groups: %v", groups)
	}

	var usersStruct []store.User
	usersByteData, _ := json.Marshal(users)
	err = json.Unmarshal(usersByteData, &usersStruct)
	if err != nil {
		logger.Fatal("Unable to parse users: %v", users)
	}

	var serversStruct []store.Server
	serversByteData, _ := json.Marshal(servers)
	err = json.Unmarshal(serversByteData, &serversStruct)
	if err != nil {
		logger.Fatal("Unable to parse servers: %v", servers)
	}

	var plansStruct []store.Plan
	plansByteData, _ := json.Marshal(plan)
	err = json.Unmarshal(plansByteData, &plansStruct)
	if err != nil {
		logger.Fatal("Unable to parse plan: %v", plan)
	}

	newStore := store.NewStore()
	newStore.
		WithUsers(usersStruct).
		WithPlans(plansStruct).
		WithServers(serversStruct).
		WithGroups(groupsStruct)

	store.Store = newStore
}

func GetAllowedCommandsForUser() map[string][]string {
	return store.Store.GetAllowedCommandsForUser(utils.GetExecutingUser())
}
