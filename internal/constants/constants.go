package constants

type UserContextKeyType string

var UserContextKey UserContextKeyType = "ctx_username"

const (
	RootStagingPath           = "/tmp/gatekeeper/"
	PrivateKeysStagingPath    = RootStagingPath + "keys/"
	PrivateKeysStagingTarPath = RootStagingPath + "keys.tar"
)
