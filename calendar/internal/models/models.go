package models

type Auth struct {
	Username string
	Password string
}

type User struct {
	Login    string
	Timezone string
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
