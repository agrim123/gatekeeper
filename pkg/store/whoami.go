package store

type Whoami struct {
	Username        string
	Groups          []string
	AllowedCommands map[string][]string
}
