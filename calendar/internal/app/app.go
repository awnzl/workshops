package app

import (
	"log"

	"calendar/internal/models"
)

type JWTGenerator interface {
	GenerateJWT(username string) (signedToken string, err error)
}

type Storage interface {
	Connect(dbURL string) error
	Close()
	LoginUser(usr, psw string) error
	User(user models.User) error
	Event(userID, eventID string) (models.Event, error)
	UpdateEvent(userID, eventID string, event models.Event) error
	AddEvents(userID string, events []models.Event) error
	Events(userID string) ([]models.Event, error)
}

type App struct {
	log  *log.Logger
	db   Storage
	auth JWTGenerator
}

func New(l *log.Logger, db Storage, a JWTGenerator) *App {
	return &App{
		log:  l,
		db:   db,
		auth: a,
	}
}

func (a *App) Login(usr, psw string) (token string, err error) {
	err = a.db.LoginUser(usr, psw)
	if err != nil {
		a.log.Println(err)
		return
	}

	return a.generateToken(usr, psw)
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

func (a *App) generateToken(usr, psw string) (token string, err error) {
	_ = psw

	token, err = a.auth.GenerateJWT(usr)
	if err != nil {
		a.log.Println(err)
		return
	}

	return
}
