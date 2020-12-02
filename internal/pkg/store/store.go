package store

type User struct {
	Email string `json:"email"`
}

type AccessMapping struct {
	User  User
	Roles []string
}

type Role struct {
	Name         string   `json:"name"`
	AllowedPlans []string `json:"allowed_plans"`
}

var Users map[string]AccessMapping
var Roles map[string]Role
var Servers map[string]*Server
var Plans map[string]Plan
