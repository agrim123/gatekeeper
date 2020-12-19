package constants

type ContextKeyType string

var UserContextKey ContextKeyType = "ctx_username"
var PlanContextKey ContextKeyType = "ctx_plan"
var OptionContextKey ContextKeyType = "ctx_option"

const (
	RootStagingPath           = "/tmp/gatekeeper/"
	PrivateKeysStagingPath    = RootStagingPath + "keys/"
	PrivateKeysStagingTarPath = RootStagingPath + "keys.tar"
	BaseImageName             = "gatekeeper"
	BaseContainerName         = "gatekeeper"
)
