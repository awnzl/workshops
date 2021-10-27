package models

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Login    string `json:"login"`
	Timezone string `json:"timezone"`
}

type Event struct {
	ID          string
	Title       string
	Description string
	Time        string
	Timezone    string
	Duration    int32
	Notes       []string
}
