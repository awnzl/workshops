package app

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"calendar/internal/auth"
	"calendar/internal/models"
)

type App struct {
	log  *log.Logger
	auth *auth.App
}

func New(l *log.Logger, a *auth.App) *App {
	return &App{
		log:  l,
		auth: a,
	}
}

// todo: replace by db
var users = make(map[string]models.User)

var tmpEvents []models.Event = []models.Event{
	{
		ID: uuid.New().String(),
		Title: "event1",
		Description: "event1 description",
		Time: time.Now().String(),
		Timezone: time.Local.String(),
		Duration: 123,
		Notes: []string{ "event1 note1", "event1 note2" },
	},
	{
		ID: uuid.New().String(),
		Title: "event2",
		Description: "event2 description",
		Time: time.Now().String(),
		Timezone: time.Local.String(),
		Duration: 15,
		Notes: []string{ "event2 note1", "event2 note2" },
	},
	{
		ID: uuid.New().String(),
		Title: "event3",
		Description: "event3 description",
		Time: time.Now().String(),
		Timezone: time.Local.String(),
		Duration: 7,
		Notes: []string{ "event3 note1", "event3 note2" },
	},
}

func (a *App) AuthApp() *auth.App {
	return a.auth
}

func (a *App) Login(usr, psw string) (token string, err error) {
	token, err = a.auth.GenerateJWT(usr)
	if err != nil {
		a.log.Println(err)
		return
	}

	users[usr] = models.User{
		Login:    usr,
		Timezone: time.Now().Local().String(),
	}

	return
}

func (a *App) Logout(usr string) error {
//todo: logout session instead deleting it
	delete(users, usr)
	return nil
}

func (a *App) User(user models.User) error {
	// tmpLoggedInUser.Timezone = user.Timezone
	//todo: update timezone in events? or just use user timezone?

	return nil
}

func (a *App) Event(id string) (event models.Event, err error) {
	for i := range tmpEvents {
		if id == tmpEvents[i].ID {
			return tmpEvents[i], nil
		}
	}

	err = fmt.Errorf("an event not found, id: %v", id)

	return
}

func (a *App) Events() ([]models.Event, error) {
	return tmpEvents, nil
}

func (a *App) AddEvents(events []models.Event) error {
	tmpEvents = append(tmpEvents, events...)

	return nil
}

func (a *App) UpdateEvent(event models.Event) error {
	for i := range tmpEvents {
		if event.ID == tmpEvents[i].ID {
			tmpEvents[i] = event
			return nil
		}
	}

	return fmt.Errorf("an event not found")
}
