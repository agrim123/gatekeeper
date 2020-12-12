package store

type user struct {
	Username string `json:"username"`
	Team     string `json:"team"`
}

type User struct {
	User   user
	Groups []string
}

type Group struct {
	Name         string   `json:"name"`
	AllowedPlans []string `json:"allowed_plans"`
}
