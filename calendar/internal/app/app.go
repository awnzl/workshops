package app

import (
	"log"

	"calendar/internal/models"
)

type Storage interface {
	Connect(dbURL string) error
	Close()
	UserToken(usr, psw string) (string, error)
	User(user models.User) error
	Event(userID, eventID string) (models.Event, error)
	UpdateEvent(userID, eventID string, event models.Event) error
	AddEvents(userID string, events []models.Event) error
	Events(userID string) ([]models.Event, error)
}

type App struct {
	log *log.Logger
	db  Storage
}

func New(l *log.Logger, db Storage) *App {
	return &App{
		log: l,
		db:  db,
	}
}

func (a *App) Login(usr, psw string) (token string, err error) {
	return a.db.UserToken(usr, psw)
}

func (a *App) Logout(userID string) error {
	return nil
}

func (a *App) User(user models.User) error {
	return a.db.User(user)
}

func (a *App) Event(userID, eventID string) (event models.Event, err error) {
	return a.db.Event(userID, eventID)
}

func (a *App) Events(userID string) ([]models.Event, error) {
	return a.db.Events(userID)
}

func (a *App) AddEvents(userID string, events []models.Event) error {
	return a.db.AddEvents(userID, events)
}

func (a *App) UpdateEvent(userID, eventID string, event models.Event) error {
	return a.db.UpdateEvent(userID, eventID, event)
}
