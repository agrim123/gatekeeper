package store

type User struct {
	Email string `json:"email"`
}

type AccessMapping struct {
	User   User
	Groups []string
}

type Group struct {
	Name         string   `json:"name"`
	AllowedPlans []string `json:"allowed_plans"`
}

var Users map[string]AccessMapping
var Groups map[string]Group
var Servers map[string]*Server
var Plans map[string]Plan
