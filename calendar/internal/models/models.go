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
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Time        string   `json:"time"`
	Timezone    string   `json:"timezone"` // depends on user's timezone
	Duration    int32    `json:"duration"`
	Notes       []string `json:"notes"`
}
